// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package hotrank

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mybilibili/app/api/creative/internal/logic/hotrank"
	"mybilibili/app/api/creative/internal/svc"
	"mybilibili/app/api/creative/internal/types"
)

// 获取全站热门排行榜
func GetHotRankListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetHotRankListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := hotrank.NewGetHotRankListLogic(r.Context(), svcCtx)
		resp, err := l.GetHotRankList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
