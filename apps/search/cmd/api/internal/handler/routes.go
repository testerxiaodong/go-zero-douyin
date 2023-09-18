// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	search "go-zero-douyin/apps/search/cmd/api/internal/handler/search"
	"go-zero-douyin/apps/search/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/video",
				Handler: search.VideoHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/video/suggestion",
				Handler: search.CompleteVideoHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/user",
				Handler: search.UserHandler(serverCtx),
			},
		},
		rest.WithPrefix("/search/v1"),
	)
}
