package search

import (
	"go-zero-douyin/common/httpResponse"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-douyin/apps/search/cmd/api/internal/logic/search"
	"go-zero-douyin/apps/search/cmd/api/internal/svc"
	"go-zero-douyin/apps/search/cmd/api/internal/types"
)

func CompleteVideoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CompleteVideoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpResponse.ParamErrorResult(r, w, err)
			return
		}

		l := search.NewCompleteVideoLogic(r.Context(), svcCtx)
		resp, err := l.CompleteVideo(&req)
		httpResponse.ApiResult(r, w, resp, err)
	}
}
