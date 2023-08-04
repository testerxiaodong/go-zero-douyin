package {{.PkgName}}

import (
    "go-zero-douyin/common/httpResponse"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	{{.ImportPackages}}
)

func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
			httpResponse.ParamErrorResult(r,w,err)
			return
		}

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
        		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
        		httpResponse.ApiResult(r ,w , resp, err)
	}
}
