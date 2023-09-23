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
				Method:  http.MethodGet,
				Path:    "/section/list",
				Handler: video.GetAllSectionHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/tag/list",
				Handler: video.GetAllTagHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/detail",
				Handler: video.DetailHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/feed",
				Handler: video.FeedHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/list",
				Handler: video.ListHandler(serverCtx),
			},
		},
		rest.WithPrefix("/video/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/section/add",
				Handler: video.AddSectionHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/section/del",
				Handler: video.DelSectionHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/tag/add",
				Handler: video.AddTagHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/tag/del",
				Handler: video.DelTagHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/publish",
				Handler: video.PublishHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/delete",
				Handler: video.DeleteVideoHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/video/v1"),
	)
}
