// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package wallet

import (
	"context"

	"mybilibili/app/wallet/cmd/api/internal/svc"
	"mybilibili/app/wallet/cmd/api/internal/types"
	walletpb "mybilibili/app/wallet/cmd/rpc/wallet"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStreamListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询流水
func NewGetStreamListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStreamListLogic {
	return &GetStreamListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStreamListLogic) GetStreamList(req *types.GetStreamListReq) (resp *types.GetStreamListResp, err error) {
	// 调用RPC服务
	rpcResp, err := l.svcCtx.WalletRpc.GetStreamList(l.ctx, &walletpb.GetStreamListReq{
		Uid:    req.Uid,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		l.Errorf("get stream list rpc failed: uid=%d, err=%v", req.Uid, err)
		return nil, err
	}

	// 转换流水列表
	list := make([]types.CoinStreamRecord, 0, len(rpcResp.List))
	for _, record := range rpcResp.List {
		list = append(list, types.CoinStreamRecord{
			Id:            record.Id,
			Uid:           record.Uid,
			TransactionId: record.TransactionId,
			CoinType:      record.CoinType,
			DeltaCoinNum:  record.DeltaCoinNum,
			OrgCoinNum:    record.OrgCoinNum,
			OpResult:      record.OpResult,
			OpReason:      record.OpReason,
			OpType:        record.OpType,
			OpTime:        record.OpTime,
			BizCode:       record.BizCode,
			Platform:      record.Platform,
		})
	}

	return &types.GetStreamListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
