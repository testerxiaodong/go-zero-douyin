package video

import (
	"go-zero-douyin/common/httpResponse"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-douyin/apps/video/cmd/api/internal/logic/video"
	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"
)

func FeedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoFeedReq
		if err := httpx.Parse(r, &req); err != nil {
			httpResponse.ParamErrorResult(r, w, err)
			return
		}

		l := video.NewFeedLogic(r.Context(), svcCtx)
		resp, err := l.Feed(&req)
		httpResponse.ApiResult(r, w, resp, err)
	}
}
