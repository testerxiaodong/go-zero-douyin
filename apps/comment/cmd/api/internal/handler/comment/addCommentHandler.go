package comment

import (
	"go-zero-douyin/common/httpResponse"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-douyin/apps/comment/cmd/api/internal/logic/comment"
	"go-zero-douyin/apps/comment/cmd/api/internal/svc"
	"go-zero-douyin/apps/comment/cmd/api/internal/types"
)

func AddCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddCommentReq
		if err := httpx.Parse(r, &req); err != nil {
			httpResponse.ParamErrorResult(r, w, err)
			return
		}

		l := comment.NewAddCommentLogic(r.Context(), svcCtx)
		err := l.AddComment(&req)
		httpResponse.ApiResult(r, w, nil, err)
	}
}
