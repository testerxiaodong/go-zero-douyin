package comment

import (
	"go-zero-douyin/common/httpResponse"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-douyin/apps/comment/cmd/api/internal/logic/comment"
	"go-zero-douyin/apps/comment/cmd/api/internal/svc"
	"go-zero-douyin/apps/comment/cmd/api/internal/types"
)

func VideoCommentListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetVideoCommentListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpResponse.ParamErrorResult(r, w, err)
			return
		}

		l := comment.NewVideoCommentListLogic(r.Context(), svcCtx)
		resp, err := l.VideoCommentList(&req)
		httpResponse.ApiResult(r, w, resp, err)
	}
}
