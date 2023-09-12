package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/common/xerr"
	"strings"

	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllVideoLogic {
	return &GetAllVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllVideoLogic) GetAllVideo(in *pb.GetAllVideoReq) (*pb.GetAllVideoResp, error) {
	// todo: add your logic here and delete this line
	// 查询数据库
	videos, err := l.svcCtx.VideoDo.GetAllVideo(l.ctx)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "查询数据库所有视频数据失败, err: %v", err)
	}

	// 没有数据，直接返回
	if len(videos) == 0 {
		return &pb.GetAllVideoResp{}, nil
	}

	// 拼接响应
	resp := &pb.GetAllVideoResp{Videos: make([]*pb.VideoInfo, 0, len(videos))}
	for _, video := range videos {
		videoInfo := &pb.VideoInfo{}
		videoInfo.Id = video.ID
		videoInfo.Title = video.Title
		videoInfo.SectionId = video.SectionID
		videoInfo.Tags = strings.Split(video.TagIds, ",")
		videoInfo.OwnerId = video.OwnerID
		videoInfo.PlayUrl = video.PlayURL
		videoInfo.CoverUrl = video.CoverURL
		videoInfo.CreateTime = video.CreateTime
		videoInfo.UpdateTime = video.UpdateTime
		resp.Videos = append(resp.Videos, videoInfo)
	}
	return resp, nil
}
