// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package search

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mybilibili/app/search/cmd/api/internal/logic/search"
	"mybilibili/app/search/cmd/api/internal/svc"
	"mybilibili/app/search/cmd/api/internal/types"
)

// 弹幕日期搜索
func DmDateSearchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DmDateSearchReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := search.NewDmDateSearchLogic(r.Context(), svcCtx)
		resp, err := l.DmDateSearch(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
