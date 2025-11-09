package logic

import (
	"context"
	"fmt"
	"time"

	"mybilibili/app/wallet/cmd/rpc/internal/svc"
	"mybilibili/app/wallet/cmd/rpc/wallet"
	"mybilibili/app/wallet/model"
	"mybilibili/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayLogic {
	return &PayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Pay 消费（严格按照主项目实现，重点：余额检查防超支）
func (l *PayLogic) Pay(in *wallet.PayReq) (*wallet.PayResp, error) {
	// 1. 参数校验
	if in.Uid <= 0 || in.CoinNum <= 0 || in.TransactionId == "" {
		l.Errorf("invalid param: uid=%d, coin_num=%d, tid=%s", in.Uid, in.CoinNum, in.TransactionId)
		return nil, xerr.NewErrCode(xerr.REQUEST_PARAM_ERROR)
	}

	// 2. 币种转换
	sysCoinType := model.GetSysCoinType(in.CoinType, in.Platform)
	sysCoinTypeNo := model.GetCoinTypeNumber(sysCoinType)
	if sysCoinTypeNo == 0 {
		l.Errorf("invalid coin_type: %s", in.CoinType)
		return nil, xerr.NewErrCode(xerr.REQUEST_PARAM_ERROR)
	}

	// 3. 锁定TransactionID（防重复提交，300秒）
	tidLockKey := fmt.Sprintf("wallet:lock:tid:%s", in.TransactionId)
	ok, err := l.svcCtx.Redis.Setnx(tidLockKey, "locked")
	if err != nil {
		l.Errorf("lock transaction id failed: tid=%s, err=%v", in.TransactionId, err)
		return nil, xerr.NewErrMsg("Redis错误")
	}
	if !ok {
		l.Errorf("transaction id already locked: tid=%s", in.TransactionId)
		return nil, xerr.NewErrMsg("交易处理中，请勿重复提交")
	}
	_ = l.svcCtx.Redis.Expire(tidLockKey, 300)
	defer l.svcCtx.Redis.Del(tidLockKey)

	// 4. 锁定用户（防并发修改，600秒）
	userLockKey := fmt.Sprintf("wallet:lock:user:%d", in.Uid)
	lockValue := fmt.Sprintf("locked:%d", time.Now().UnixNano())
	ok, err = l.svcCtx.Redis.Setnx(userLockKey, lockValue)
	if err != nil || !ok {
		l.Errorf("lock user failed: uid=%d, err=%v", in.Uid, err)
		return nil, xerr.NewErrMsg("账户操作中，请稍后重试")
	}
	_ = l.svcCtx.Redis.Expire(userLockKey, 600)
	defer func() {
		val, _ := l.svcCtx.Redis.Get(userLockKey)
		if val == lockValue {
			l.svcCtx.Redis.Del(userLockKey)
		}
	}()

	// 5. 准备流水记录（默认失败）
	stream := &model.CoinStreamRecord{
		Uid:           in.Uid,
		TransactionId: in.TransactionId,
		ExtendTid:     in.ExtendTid,
		CoinType:      sysCoinTypeNo,
		DeltaCoinNum:  -in.CoinNum, // 负数表示扣减
		OpType:        model.OpTypePay,
		OpTime:        time.Now(),
		BizCode:       in.BizCode,
		Platform:      model.GetPlatformNumber(in.Platform),
		OpResult:      model.OpResultSubFailed, // 默认失败
		OpReason:      model.OpReasonSuccess,
	}

	// 6. 执行数据库事务
	var walletResp *model.UserWallet
	err = l.svcCtx.UserWalletModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 6.1 FOR UPDATE查询（悲观锁）
		walletData, err := l.svcCtx.UserWalletModel.FindOneForUpdate(ctx, session, in.Uid)
		if err == model.ErrNotFound {
			// 用户钱包不存在
			stream.OpReason = model.OpReasonInvalidParam
			l.Errorf("wallet not found: uid=%d", in.Uid)
			// 记录失败流水
			_, _ = l.svcCtx.CoinStreamRecordModel.Insert(ctx, session, stream)
			return xerr.NewErrCode(xerr.WALLET_NOT_FOUND)
		} else if err != nil {
			l.Errorf("find wallet failed: uid=%d, err=%v", in.Uid, err)
			return err
		}

		// 记录原余额
		stream.OrgCoinNum = model.GetCoinByType(walletData, sysCoinTypeNo)

		// 6.2 余额检查（防超支）⭐关键
		if stream.OrgCoinNum < in.CoinNum {
			stream.OpReason = model.OpReasonNotEnough
			l.Errorf("coin not enough: uid=%d, coin_type=%d, balance=%d, need=%d",
				in.Uid, sysCoinTypeNo, stream.OrgCoinNum, in.CoinNum)
			// 记录失败流水
			_, _ = l.svcCtx.CoinStreamRecordModel.Insert(ctx, session, stream)
			return xerr.NewErrCode(xerr.COIN_NOT_ENOUGH)
		}

		// 6.3 UPDATE扣款
		err = l.svcCtx.UserWalletModel.UpdatePay(ctx, session, in.Uid, sysCoinTypeNo, in.CoinNum)
		if err != nil {
			l.Errorf("update pay failed: uid=%d, coin_type=%d, amount=%d, err=%v",
				in.Uid, sysCoinTypeNo, in.CoinNum, err)
			return err
		}

		// 6.4 INSERT成功流水
		stream.OpResult = model.OpResultSubSucc
		_, err = l.svcCtx.CoinStreamRecordModel.Insert(ctx, session, stream)
		if err != nil {
			l.Errorf("insert stream failed: uid=%d, err=%v", in.Uid, err)
			return err
		}

		// 更新内存中的wallet（用于返回）
		model.SubCoin(walletData, sysCoinTypeNo, in.CoinNum)
		walletResp = walletData
		return nil
	})

	if err != nil {
		// 判断是否是业务错误（余额不足、钱包不存在）
		if e, ok := err.(*xerr.CodeError); ok {
			if e.GetErrCode() == xerr.COIN_NOT_ENOUGH || e.GetErrCode() == xerr.WALLET_NOT_FOUND {
				return nil, err // 直接返回业务错误
			}
		}
		l.Errorf("pay transaction failed: uid=%d, err=%v", in.Uid, err)
		// 异步记录失败流水
		go func() {
			stream.OpResult = model.OpResultSubFailed
			stream.OpReason = model.OpReasonInvalidParam
			_, _ = l.svcCtx.CoinStreamRecordModel.Insert(context.Background(), nil, stream)
		}()
		return nil, xerr.NewErrCode(xerr.DB_ERROR)
	}

	// 7. 删除缓存（如果使用缓存）
	// l.svcCtx.WalletCache.Del(l.ctx, in.Uid)

	// 8. 发布钱包变更消息（TODO: 实现Kafka发布）
	// go l.publishWalletChange(in.Uid, "pay", in.CoinType, in.CoinNum, walletResp)

	// 9. 返回新余额
	l.Infof("pay success: uid=%d, coin_type=%s, amount=%d, new_balance: gold=%d, iap_gold=%d, silver=%d",
		in.Uid, sysCoinType, in.CoinNum, walletResp.Gold, walletResp.IapGold, walletResp.Silver)

	return &wallet.PayResp{
		Gold:    walletResp.Gold,
		IapGold: walletResp.IapGold,
		Silver:  walletResp.Silver,
	}, nil
}
