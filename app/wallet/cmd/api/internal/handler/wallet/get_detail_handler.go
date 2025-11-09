// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package wallet

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mybilibili/app/wallet/cmd/api/internal/logic/wallet"
	"mybilibili/app/wallet/cmd/api/internal/svc"
	"mybilibili/app/wallet/cmd/api/internal/types"
)

// 查询余额
func GetDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := wallet.NewGetDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
