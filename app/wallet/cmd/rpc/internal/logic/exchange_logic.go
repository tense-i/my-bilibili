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

type ExchangeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExchangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeLogic {
	return &ExchangeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Exchange 兑换（严格按照主项目实现，原子性操作）
func (l *ExchangeLogic) Exchange(in *wallet.ExchangeReq) (*wallet.ExchangeResp, error) {
	// 1. 参数校验
	if in.Uid <= 0 || in.SrcCoinNum <= 0 || in.DestCoinNum <= 0 || in.TransactionId == "" {
		l.Errorf("invalid param: uid=%d, src_num=%d, dest_num=%d, tid=%s",
			in.Uid, in.SrcCoinNum, in.DestCoinNum, in.TransactionId)
		return nil, xerr.NewErrCode(xerr.REQUEST_PARAM_ERROR)
	}

	// 2. 币种转换
	srcSysCoinType := model.GetSysCoinType(in.SrcCoinType, in.Platform)
	srcCoinTypeNo := model.GetCoinTypeNumber(srcSysCoinType)
	destSysCoinType := model.GetSysCoinType(in.DestCoinType, in.Platform)
	destCoinTypeNo := model.GetCoinTypeNumber(destSysCoinType)

	if srcCoinTypeNo == 0 || destCoinTypeNo == 0 {
		l.Errorf("invalid coin_type: src=%s, dest=%s", in.SrcCoinType, in.DestCoinType)
		return nil, xerr.NewErrCode(xerr.INVALID_COIN_TYPE)
	}

	// 3. 禁止相同币种兑换
	if srcCoinTypeNo == destCoinTypeNo {
		l.Errorf("cannot exchange same coin type: %s", srcSysCoinType)
		return nil, xerr.NewErrCode(xerr.INVALID_COIN_TYPE)
	}

	// 4. 验证兑换比例（通常是1:1）
	if in.SrcCoinNum != in.DestCoinNum {
		l.Errorf("invalid exchange rate: src=%d, dest=%d", in.SrcCoinNum, in.DestCoinNum)
		return nil, xerr.NewErrCode(xerr.EXCHANGE_RATE_ERROR)
	}

	// 5. 锁定TransactionID（防重复提交，300秒）
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

	// 6. 锁定用户（防并发修改，600秒）
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

	// 7. 准备流水记录（源币种扣除和目标币种增加）
	srcStream := &model.CoinStreamRecord{
		Uid:           in.Uid,
		TransactionId: in.TransactionId,
		CoinType:      srcCoinTypeNo,
		DeltaCoinNum:  -in.SrcCoinNum, // 负数表示扣减
		OpType:        model.OpTypeExchange,
		OpTime:        time.Now(),
		Platform:      model.GetPlatformNumber(in.Platform),
		OpResult:      model.OpResultSubFailed, // 默认失败
		OpReason:      model.OpReasonSuccess,
	}

	destStream := &model.CoinStreamRecord{
		Uid:           in.Uid,
		TransactionId: in.TransactionId,
		CoinType:      destCoinTypeNo,
		DeltaCoinNum:  in.DestCoinNum, // 正数表示增加
		OpType:        model.OpTypeExchange,
		OpTime:        time.Now(),
		Platform:      model.GetPlatformNumber(in.Platform),
		OpResult:      model.OpResultAddFailed, // 默认失败
		OpReason:      model.OpReasonSuccess,
	}

	// 8. 准备兑换记录
	exchangeRecord := &model.CoinExchangeRecord{
		Uid:           in.Uid,
		TransactionId: in.TransactionId,
		SrcCoinType:   srcCoinTypeNo,
		SrcCoinNum:    in.SrcCoinNum,
		DestCoinType:  destCoinTypeNo,
		DestCoinNum:   in.DestCoinNum,
		ExchangeRate:  1.0,
		Status:        0, // 默认失败
	}

	// 9. 执行数据库事务（原子性）
	var walletResp *model.UserWallet
	err = l.svcCtx.UserWalletModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 9.1 FOR UPDATE查询
		walletData, err := l.svcCtx.UserWalletModel.FindOneForUpdate(ctx, session, in.Uid)
		if err == model.ErrNotFound {
			srcStream.OpReason = model.OpReasonInvalidParam
			destStream.OpReason = model.OpReasonInvalidParam
			l.Errorf("wallet not found: uid=%d", in.Uid)
			return xerr.NewErrCode(xerr.WALLET_NOT_FOUND)
		} else if err != nil {
			l.Errorf("find wallet failed: uid=%d, err=%v", in.Uid, err)
			return err
		}

		// 记录原余额
		srcStream.OrgCoinNum = model.GetCoinByType(walletData, srcCoinTypeNo)
		destStream.OrgCoinNum = model.GetCoinByType(walletData, destCoinTypeNo)

		// 9.2 源币种余额检查
		if srcStream.OrgCoinNum < in.SrcCoinNum {
			srcStream.OpReason = model.OpReasonNotEnough
			destStream.OpReason = model.OpReasonNotEnough
			l.Errorf("coin not enough: uid=%d, coin_type=%d, balance=%d, need=%d",
				in.Uid, srcCoinTypeNo, srcStream.OrgCoinNum, in.SrcCoinNum)
			// 记录失败流水
			_, _ = l.svcCtx.CoinStreamRecordModel.Insert(ctx, session, srcStream)
			_, _ = l.svcCtx.CoinStreamRecordModel.Insert(ctx, session, destStream)
			return xerr.NewErrCode(xerr.COIN_NOT_ENOUGH)
		}

		// 9.3 UPDATE兑换（原子性：一条SQL完成）
		err = l.svcCtx.UserWalletModel.UpdateExchange(ctx, session, in.Uid,
			srcCoinTypeNo, in.SrcCoinNum, destCoinTypeNo, in.DestCoinNum)
		if err != nil {
			l.Errorf("update exchange failed: uid=%d, err=%v", in.Uid, err)
			return err
		}

		// 9.4 INSERT成功流水（两条）
		srcStream.OpResult = model.OpResultSubSucc
		destStream.OpResult = model.OpResultAddSucc
		_, err = l.svcCtx.CoinStreamRecordModel.Insert(ctx, session, srcStream)
		if err != nil {
			l.Errorf("insert src stream failed: uid=%d, err=%v", in.Uid, err)
			return err
		}
		_, err = l.svcCtx.CoinStreamRecordModel.Insert(ctx, session, destStream)
		if err != nil {
			l.Errorf("insert dest stream failed: uid=%d, err=%v", in.Uid, err)
			return err
		}

		// 9.5 INSERT兑换记录
		exchangeRecord.Status = 1 // 成功
		_, err = l.svcCtx.CoinExchangeRecordModel.Insert(ctx, session, exchangeRecord)
		if err != nil {
			l.Errorf("insert exchange record failed: uid=%d, err=%v", in.Uid, err)
			return err
		}

		// 更新内存中的wallet（用于返回）
		model.SubCoin(walletData, srcCoinTypeNo, in.SrcCoinNum)
		model.AddCoin(walletData, destCoinTypeNo, in.DestCoinNum)
		walletResp = walletData
		return nil
	})

	if err != nil {
		// 判断是否是业务错误
		if e, ok := err.(*xerr.CodeError); ok {
			if e.GetErrCode() == xerr.COIN_NOT_ENOUGH || e.GetErrCode() == xerr.WALLET_NOT_FOUND {
				return nil, err
			}
		}
		l.Errorf("exchange transaction failed: uid=%d, err=%v", in.Uid, err)
		return nil, xerr.NewErrCode(xerr.DB_ERROR)
	}

	// 10. 返回新余额
	l.Infof("exchange success: uid=%d, %s(%d) -> %s(%d), new_balance: gold=%d, iap_gold=%d, silver=%d",
		in.Uid, srcSysCoinType, in.SrcCoinNum, destSysCoinType, in.DestCoinNum,
		walletResp.Gold, walletResp.IapGold, walletResp.Silver)

	return &wallet.ExchangeResp{
		Gold:    walletResp.Gold,
		IapGold: walletResp.IapGold,
		Silver:  walletResp.Silver,
	}, nil
}
