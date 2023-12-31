// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package server

import (
	"context"

	"go-zero-douyin/apps/user/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/user/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServer) RegisterUser(ctx context.Context, in *pb.RegisterUserReq) (*pb.RegisterUserResp, error) {
	l := logic.NewRegisterUserLogic(ctx, s.svcCtx)
	return l.RegisterUser(in)
}

func (s *UserServer) GenerateToken(ctx context.Context, in *pb.GenerateTokenReq) (*pb.GenerateTokenResp, error) {
	l := logic.NewGenerateTokenLogic(ctx, s.svcCtx)
	return l.GenerateToken(in)
}

func (s *UserServer) Login(ctx context.Context, in *pb.LoginReq) (*pb.LoginResp, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}

func (s *UserServer) GetUserInfo(ctx context.Context, in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	l := logic.NewGetUserInfoLogic(ctx, s.svcCtx)
	return l.GetUserInfo(in)
}

func (s *UserServer) UpdateUser(ctx context.Context, in *pb.UpdateUserReq) (*pb.UpdateUserResp, error) {
	l := logic.NewUpdateUserLogic(ctx, s.svcCtx)
	return l.UpdateUser(in)
}
