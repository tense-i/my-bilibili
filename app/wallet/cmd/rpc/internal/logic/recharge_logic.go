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

type RechargeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRechargeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RechargeLogic {
	return &RechargeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Recharge 充值（严格按照主项目实现）
func (l *RechargeLogic) Recharge(in *wallet.RechargeReq) (*wallet.RechargeResp, error) {
	// 1. 参数校验
	if in.Uid <= 0 || in.CoinNum <= 0 || in.TransactionId == "" {
		l.Errorf("invalid param: uid=%d, coin_num=%d, tid=%s", in.Uid, in.CoinNum, in.TransactionId)
		return nil, xerr.NewErrCode(xerr.REQUEST_PARAM_ERROR)
	}

	// 2. 币种转换（iOS平台的gold自动转为iap_gold）
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

	// 5. 准备流水记录
	stream := &model.CoinStreamRecord{
		Uid:           in.Uid,
		TransactionId: in.TransactionId,
		ExtendTid:     in.ExtendTid,
		CoinType:      sysCoinTypeNo,
		DeltaCoinNum:  in.CoinNum,
		OpType:        model.OpTypeRecharge,
		OpTime:        time.Now(),
		BizCode:       in.BizCode,
		Platform:      model.GetPlatformNumber(in.Platform),
		OpResult:      model.OpResultAddFailed, // 默认失败
		OpReason:      model.OpReasonSuccess,
	}

	// 6. 执行数据库事务
	var walletResp *model.UserWallet
	err = l.svcCtx.UserWalletModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 6.1 FOR UPDATE查询（悲观锁）
		walletData, err := l.svcCtx.UserWalletModel.FindOneForUpdate(ctx, session, in.Uid)
		if err == model.ErrNotFound {
			// 首次充值，初始化钱包
			walletData = &model.UserWallet{
				Uid:     in.Uid,
				Gold:    0,
				IapGold: 0,
				Silver:  0,
			}
			_, err = l.svcCtx.UserWalletModel.Insert(ctx, session, walletData)
			if err != nil {
				l.Errorf("insert wallet failed: uid=%d, err=%v", in.Uid, err)
				return err
			}
			// 重新查询
			walletData, err = l.svcCtx.UserWalletModel.FindOneForUpdate(ctx, session, in.Uid)
			if err != nil {
				return err
			}
		} else if err != nil {
			l.Errorf("find wallet failed: uid=%d, err=%v", in.Uid, err)
			return err
		}

		// 记录原余额
		stream.OrgCoinNum = model.GetCoinByType(walletData, sysCoinTypeNo)

		// 6.2 UPDATE余额
		err = l.svcCtx.UserWalletModel.UpdateRecharge(ctx, session, in.Uid, sysCoinTypeNo, in.CoinNum)
		if err != nil {
			l.Errorf("update recharge failed: uid=%d, coin_type=%d, amount=%d, err=%v",
				in.Uid, sysCoinTypeNo, in.CoinNum, err)
			return err
		}

		// 6.3 INSERT流水记录（成功）
		stream.OpResult = model.OpResultAddSucc
		_, err = l.svcCtx.CoinStreamRecordModel.Insert(ctx, session, stream)
		if err != nil {
			l.Errorf("insert stream failed: uid=%d, err=%v", in.Uid, err)
			return err
		}

		// 更新内存中的wallet（用于返回）
		model.AddCoin(walletData, sysCoinTypeNo, in.CoinNum)
		walletResp = walletData
		return nil
	})

	if err != nil {
		l.Errorf("recharge transaction failed: uid=%d, err=%v", in.Uid, err)
		// 异步记录失败流水
		go func() {
			stream.OpResult = model.OpResultAddFailed
			stream.OpReason = model.OpReasonInvalidParam
			_, _ = l.svcCtx.CoinStreamRecordModel.Insert(context.Background(), nil, stream)
		}()
		return nil, xerr.NewErrCode(xerr.DB_ERROR)
	}

	// 7. 删除缓存（如果使用缓存）
	// l.svcCtx.WalletCache.Del(l.ctx, in.Uid)

	// 8. 发布钱包变更消息（TODO: 实现Kafka发布）
	// go l.publishWalletChange(in.Uid, "recharge", in.CoinType, in.CoinNum, walletResp)

	// 9. 返回新余额
	l.Infof("recharge success: uid=%d, coin_type=%s, amount=%d, new_balance: gold=%d, iap_gold=%d, silver=%d",
		in.Uid, sysCoinType, in.CoinNum, walletResp.Gold, walletResp.IapGold, walletResp.Silver)

	return &wallet.RechargeResp{
		Gold:    walletResp.Gold,
		IapGold: walletResp.IapGold,
		Silver:  walletResp.Silver,
	}, nil
}
