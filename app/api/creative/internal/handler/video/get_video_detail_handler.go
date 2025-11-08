// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mybilibili/app/api/creative/internal/logic/video"
	"mybilibili/app/api/creative/internal/svc"
	"mybilibili/app/api/creative/internal/types"
)

// 获取视频详情
func GetVideoDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetVideoDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := video.NewGetVideoDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetVideoDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
