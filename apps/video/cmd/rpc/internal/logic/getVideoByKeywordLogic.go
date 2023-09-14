package logic

import (
	"context"

	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoByKeywordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoByKeywordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoByKeywordLogic {
	return &GetVideoByKeywordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoByKeywordLogic) GetVideoByKeyword(in *pb.GetVideoByKeywordReq) (*pb.GetVideoByKeywordResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetVideoByKeywordResp{}, nil
}
