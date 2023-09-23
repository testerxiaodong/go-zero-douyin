package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoByIdLogic {
	return &GetVideoByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoByIdLogic) GetVideoById(in *pb.GetVideoByIdReq) (*pb.GetVideoByIdResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get video by id with empty param")
	}
	if in.GetId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get video by id with empty id")
	}

	// 查询数据库
	video, err := l.svcCtx.VideoModel.FindOne(l.ctx, in.GetId())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "find video by id failed: %v", err)
	}

	// 视频不存在
	if video == nil {
		return nil, errors.Wrapf(ErrVideoNotFound, "video_id: %d", in.GetId())
	}

	// 返回数据
	resp := &pb.GetVideoByIdResp{Video: &pb.VideoInfo{}}
	resp.Video.Id = video.Id
	resp.Video.Title = video.Title
	resp.Video.SectionId = video.SectionId
	resp.Video.TagIds = video.TagIds
	resp.Video.OwnerId = video.OwnerId
	resp.Video.OwnerName = video.OwnerName
	resp.Video.PlayUrl = video.PlayUrl
	resp.Video.CoverUrl = video.CoverUrl
	resp.Video.CreateTime = video.CreateTime.Unix()
	resp.Video.UpdateTime = video.UpdateTime.Unix()
	return resp, nil
}
