package video

import (
	"go-zero-douyin/common/httpResponse"
	"net/http"

	"go-zero-douyin/apps/video/cmd/api/internal/logic/video"
	"go-zero-douyin/apps/video/cmd/api/internal/svc"
)

func GetAllTagHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := video.NewGetAllTagLogic(r.Context(), svcCtx)
		resp, err := l.GetAllTag()
		httpResponse.ApiResult(r, w, resp, err)
	}
}
