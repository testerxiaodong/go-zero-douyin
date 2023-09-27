package recommend

import (
	"go-zero-douyin/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-douyin/apps/recommend/cmd/api/internal/logic/recommend"
	"go-zero-douyin/apps/recommend/cmd/api/internal/svc"
	"go-zero-douyin/apps/recommend/cmd/api/internal/types"
)

func VideoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoRecommendReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := recommend.NewVideoLogic(r.Context(), svcCtx)
		resp, err := l.Video(&req)
		response.ApiResult(r, w, resp, err)
	}
}
