package video

import (
	"go-zero-douyin/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-douyin/apps/video/cmd/api/internal/logic/video"
	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"
)

func AddTagHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddTagReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := video.NewAddTagLogic(r.Context(), svcCtx)
		err := l.AddTag(&req)
		response.ApiResult(r, w, nil, err)
	}
}
