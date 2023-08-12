// Code generated by goctl. DO NOT EDIT.
// Source: videorpc.proto

package server

import (
	"context"

	"go-zero-douyin/apps/video/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
)

type VideorpcServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedVideorpcServer
}

func NewVideorpcServer(svcCtx *svc.ServiceContext) *VideorpcServer {
	return &VideorpcServer{
		svcCtx: svcCtx,
	}
}

func (s *VideorpcServer) PublishVideo(ctx context.Context, in *pb.PublishVideoReq) (*pb.PublishVideoResp, error) {
	l := logic.NewPublishVideoLogic(ctx, s.svcCtx)
	return l.PublishVideo(in)
}

func (s *VideorpcServer) VideoFeed(ctx context.Context, in *pb.VideoFeedReq) (*pb.VideoFeedResp, error) {
	l := logic.NewVideoFeedLogic(ctx, s.svcCtx)
	return l.VideoFeed(in)
}

func (s *VideorpcServer) UserVideoList(ctx context.Context, in *pb.UserVideoListReq) (*pb.UserVideoListResp, error) {
	l := logic.NewUserVideoListLogic(ctx, s.svcCtx)
	return l.UserVideoList(in)
}
