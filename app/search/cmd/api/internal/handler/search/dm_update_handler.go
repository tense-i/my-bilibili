package search

import (
	"net/http"

	"mybilibili/app/search/cmd/api/internal/logic/search"
	"mybilibili/app/search/cmd/api/internal/svc"
	"mybilibili/app/search/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DmUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DmUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := search.NewDmUpdateLogic(r.Context(), svcCtx)
		resp, err := l.DmUpdate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
