package follow

import (
	"go-zero-douyin/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-douyin/apps/social/cmd/api/internal/logic/follow"
	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"
)

func GetUserFollowCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserFollowCountReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := follow.NewGetUserFollowCountLogic(r.Context(), svcCtx)
		resp, err := l.GetUserFollowCount(&req)
		response.ApiResult(r, w, resp, err)
	}
}
