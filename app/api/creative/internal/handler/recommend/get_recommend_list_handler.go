// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package recommend

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mybilibili/app/api/creative/internal/logic/recommend"
	"mybilibili/app/api/creative/internal/svc"
	"mybilibili/app/api/creative/internal/types"
)

// 获取推荐列表
func GetRecommendListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRecommendListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := recommend.NewGetRecommendListLogic(r.Context(), svcCtx)
		resp, err := l.GetRecommendList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
