package video

import (
	"go-zero-douyin/common/httpResponse"
	"net/http"

	"go-zero-douyin/apps/video/cmd/api/internal/logic/video"
	"go-zero-douyin/apps/video/cmd/api/internal/svc"
)

func GetAllSectionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := video.NewGetAllSectionLogic(r.Context(), svcCtx)
		resp, err := l.GetAllSection()
		httpResponse.ApiResult(r, w, resp, err)
	}
}
