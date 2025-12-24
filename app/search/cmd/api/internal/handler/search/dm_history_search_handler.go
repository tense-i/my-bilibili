package search

import (
	"net/http"

	"mybilibili/app/search/cmd/api/internal/logic/search"
	"mybilibili/app/search/cmd/api/internal/svc"
	"mybilibili/app/search/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DmHistorySearchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DmHistorySearchReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := search.NewDmHistorySearchLogic(r.Context(), svcCtx)
		resp, err := l.DmHistorySearch(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
