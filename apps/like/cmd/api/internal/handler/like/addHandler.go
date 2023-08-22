package like

import (
	"go-zero-douyin/common/httpResponse"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-douyin/apps/like/cmd/api/internal/logic/like"
	"go-zero-douyin/apps/like/cmd/api/internal/svc"
	"go-zero-douyin/apps/like/cmd/api/internal/types"
)

func AddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoLikeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpResponse.ParamErrorResult(r, w, err)
			return
		}

		l := like.NewAddLogic(r.Context(), svcCtx)
		err := l.Add(&req)
		httpResponse.ApiResult(r, w, nil, err)
	}
}
