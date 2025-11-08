package logic

import (
	"context"

	"mybilibili/app/hotrank/cmd/rpc/hotrank"
	"mybilibili/app/hotrank/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHotByOIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetHotByOIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotByOIDLogic {
	return &GetHotByOIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取热度值（根据OID）
func (l *GetHotByOIDLogic) GetHotByOID(in *hotrank.GetHotByOIDReq) (*hotrank.GetHotByOIDResp, error) {
	// 查询热度值
	hot, err := l.svcCtx.AcademyArchiveModel.FindHotByOID(l.ctx, in.Oid, int(in.Business))
	if err != nil {
		l.Errorf("GetHotByOID FindHotByOID error: %v", err)
		return nil, err
	}

	return &hotrank.GetHotByOIDResp{Hot: hot}, nil
}
