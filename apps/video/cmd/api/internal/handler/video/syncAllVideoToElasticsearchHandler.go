package video

import (
	"go-zero-douyin/common/httpResponse"
	"net/http"

	"go-zero-douyin/apps/video/cmd/api/internal/logic/video"
	"go-zero-douyin/apps/video/cmd/api/internal/svc"
)

func SyncAllVideoToElasticsearchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := video.NewSyncAllVideoToElasticsearchLogic(r.Context(), svcCtx)
		err := l.SyncAllVideoToElasticsearch()
		httpResponse.ApiResult(r, w, nil, err)
	}
}
