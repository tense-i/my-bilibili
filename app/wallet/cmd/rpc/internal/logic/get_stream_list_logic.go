package logic

import (
	"context"

	"mybilibili/app/wallet/cmd/rpc/internal/svc"
	"mybilibili/app/wallet/cmd/rpc/wallet"
	"mybilibili/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStreamListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStreamListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStreamListLogic {
	return &GetStreamListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetStreamList 查询流水列表（支持分页）
func (l *GetStreamListLogic) GetStreamList(in *wallet.GetStreamListReq) (*wallet.GetStreamListResp, error) {
	// 1. 参数校验
	if in.Uid <= 0 {
		l.Errorf("invalid param: uid=%d", in.Uid)
		return nil, xerr.NewErrCode(xerr.REQUEST_PARAM_ERROR)
	}

	// 2. 设置默认分页参数
	offset := in.Offset
	limit := in.Limit
	if limit <= 0 || limit > 100 {
		limit = 20 // 默认20条，最多100条
	}

	// 3. 查询流水记录（需要查询所有分表）
	records, err := l.svcCtx.CoinStreamRecordModel.FindByUid(l.ctx, in.Uid, offset, limit)
	if err != nil {
		l.Errorf("find stream records failed: uid=%d, err=%v", in.Uid, err)
		return nil, xerr.NewErrCode(xerr.DB_ERROR)
	}

	// 4. 查询总数
	total, err := l.svcCtx.CoinStreamRecordModel.Count(l.ctx, in.Uid)
	if err != nil {
		l.Errorf("count stream records failed: uid=%d, err=%v", in.Uid, err)
		return nil, xerr.NewErrCode(xerr.DB_ERROR)
	}

	// 5. 转换为protobuf格式
	list := make([]*wallet.CoinStreamRecord, 0, len(records))
	for _, record := range records {
		list = append(list, &wallet.CoinStreamRecord{
			Id:            record.Id,
			Uid:           record.Uid,
			TransactionId: record.TransactionId,
			CoinType:      record.CoinType,
			DeltaCoinNum:  record.DeltaCoinNum,
			OrgCoinNum:    record.OrgCoinNum,
			OpResult:      record.OpResult,
			OpReason:      record.OpReason,
			OpType:        record.OpType,
			OpTime:        record.OpTime.Format("2006-01-02 15:04:05"),
			BizCode:       record.BizCode,
			Platform:      record.Platform,
		})
	}

	l.Infof("get stream list success: uid=%d, count=%d, total=%d", in.Uid, len(list), total)

	return &wallet.GetStreamListResp{
		List:  list,
		Total: total,
	}, nil
}
