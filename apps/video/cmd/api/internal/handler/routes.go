// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	video "go-zero-douyin/apps/video/cmd/api/internal/handler/video"
	"go-zero-douyin/apps/video/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/publish",
				Handler: video.PublishHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/video/v1"),
	)
}