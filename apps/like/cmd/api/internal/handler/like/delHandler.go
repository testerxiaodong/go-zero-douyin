package like

import (
	"go-zero-douyin/common/httpResponse"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-douyin/apps/like/cmd/api/internal/logic/like"
	"go-zero-douyin/apps/like/cmd/api/internal/svc"
	"go-zero-douyin/apps/like/cmd/api/internal/types"
)

func DelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoUnlikeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpResponse.ParamErrorResult(r, w, err)
			return
		}

		l := like.NewDelLogic(r.Context(), svcCtx)
		err := l.Del(&req)
		httpResponse.ApiResult(r, w, nil, err)
	}
}
