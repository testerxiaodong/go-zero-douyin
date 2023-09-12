package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/common/xerr"
	"strings"

	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishVideoLogic {
	return &PublishVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishVideoLogic) PublishVideo(in *pb.PublishVideoReq) (*pb.PublishVideoResp, error) {
	// todo: add your logic here and delete this line
	// 校验输入
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Publish video empty param")
	}

	// 插入数据
	video := &model.Video{
		Title:     in.GetTitle(),
		OwnerID:   in.GetOwnerId(),
		SectionID: in.GetSectionId(),
		TagIds:    strings.Join(in.GetTags(), ","),
		PlayURL:   in.GetPlayUrl(),
		CoverURL:  in.GetCoverUrl(),
	}
	err := l.svcCtx.VideoDo.InsertVideo(l.ctx, video)

	// 插入失败
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "insert video failed, err: %v", err)
	}

	// 返回信息
	return &pb.PublishVideoResp{
		Video: &pb.VideoInfo{
			Id:        video.ID,
			Title:     video.Title,
			SectionId: video.SectionID,
			Tags:      strings.Split(video.TagIds, ","),
			OwnerId:   video.OwnerID,
			PlayUrl:   video.PlayURL,
			CoverUrl:  video.CoverURL,
		},
	}, nil
}
