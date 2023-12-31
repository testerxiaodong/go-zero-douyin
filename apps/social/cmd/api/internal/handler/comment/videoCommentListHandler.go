package comment

import (
	"go-zero-douyin/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-douyin/apps/social/cmd/api/internal/logic/comment"
	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"
)

func VideoCommentListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetVideoCommentListReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := comment.NewVideoCommentListLogic(r.Context(), svcCtx)
		resp, err := l.VideoCommentList(&req)
		response.ApiResult(r, w, resp, err)
	}
}
