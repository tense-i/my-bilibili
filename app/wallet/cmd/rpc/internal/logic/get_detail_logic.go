package logic

import (
	"context"

	"mybilibili/app/wallet/cmd/rpc/internal/svc"
	"mybilibili/app/wallet/cmd/rpc/wallet"
	"mybilibili/app/wallet/model"
	"mybilibili/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDetailLogic {
	return &GetDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetDetail 查询余额详情（只读操作，无需加锁）
func (l *GetDetailLogic) GetDetail(in *wallet.GetDetailReq) (*wallet.GetDetailResp, error) {
	// 1. 参数校验
	if in.Uid <= 0 {
		l.Errorf("invalid param: uid=%d", in.Uid)
		return nil, xerr.NewErrCode(xerr.REQUEST_PARAM_ERROR)
	}

	// 2. 查询钱包详情（先尝试从缓存读取）
	walletData, err := l.svcCtx.UserWalletModel.FindOne(l.ctx, in.Uid)
	if err == model.ErrNotFound {
		l.Infof("wallet not found, return zero balance: uid=%d", in.Uid)
		// 钱包不存在，返回零余额（避免报错，提升用户体验）
		return &wallet.GetDetailResp{
			Detail: &wallet.WalletDetail{
				Uid:             in.Uid,
				Gold:            0,
				IapGold:         0,
				Silver:          0,
				GoldRechargeCnt: 0,
				GoldPayCnt:      0,
				SilverPayCnt:    0,
				CostBase:        0,
			},
		}, nil
	} else if err != nil {
		l.Errorf("find wallet failed: uid=%d, err=%v", in.Uid, err)
		return nil, xerr.NewErrCode(xerr.DB_ERROR)
	}

	// 3. 返回详情
	l.Infof("get detail success: uid=%d, gold=%d, iap_gold=%d, silver=%d",
		in.Uid, walletData.Gold, walletData.IapGold, walletData.Silver)

	return &wallet.GetDetailResp{
		Detail: &wallet.WalletDetail{
			Uid:             walletData.Uid,
			Gold:            walletData.Gold,
			IapGold:         walletData.IapGold,
			Silver:          walletData.Silver,
			GoldRechargeCnt: walletData.GoldRechargeCnt,
			GoldPayCnt:      walletData.GoldPayCnt,
			SilverPayCnt:    walletData.SilverPayCnt,
			CostBase:        walletData.CostBase,
		},
	}, nil
}
