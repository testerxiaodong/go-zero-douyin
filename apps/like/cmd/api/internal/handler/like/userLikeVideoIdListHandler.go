package like

import (
	"go-zero-douyin/common/httpResponse"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-douyin/apps/like/cmd/api/internal/logic/like"
	"go-zero-douyin/apps/like/cmd/api/internal/svc"
	"go-zero-douyin/apps/like/cmd/api/internal/types"
)

func UserLikeVideoIdListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserLikeVideoIdListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpResponse.ParamErrorResult(r, w, err)
			return
		}

		l := like.NewUserLikeVideoIdListLogic(r.Context(), svcCtx)
		resp, err := l.UserLikeVideoIdList(&req)
		httpResponse.ApiResult(r, w, resp, err)
	}
}
