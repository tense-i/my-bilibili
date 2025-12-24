// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package allowance

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mybilibili/app/coupon/cmd/api/internal/logic/allowance"
	"mybilibili/app/coupon/cmd/api/internal/svc"
	"mybilibili/app/coupon/cmd/api/internal/types"
)

// 领取代金券
func ReceiveAllowanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReceiveAllowanceReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := allowance.NewReceiveAllowanceLogic(r.Context(), svcCtx)
		resp, err := l.ReceiveAllowance(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
