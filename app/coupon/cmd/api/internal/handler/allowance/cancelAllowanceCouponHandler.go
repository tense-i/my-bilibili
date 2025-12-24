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

// 取消使用代金券
func CancelAllowanceCouponHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CancelAllowanceCouponReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := allowance.NewCancelAllowanceCouponLogic(r.Context(), svcCtx)
		err := l.CancelAllowanceCoupon(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
