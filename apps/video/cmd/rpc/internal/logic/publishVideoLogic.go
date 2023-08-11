package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/common/xerr"

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
	video := &model.Video{}
	videoQuery := l.svcCtx.Query.Video
	err := copier.Copy(video, in)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "复制pb 到 db失败，err: %v", err)
	}
	err = videoQuery.WithContext(l.ctx).Create(video)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "insert video failed, err: %v", err)
	}
	// 返回信息
	return &pb.PublishVideoResp{
		Video: &pb.Video{
			Id:       video.ID,
			Title:    video.Title,
			OwnerId:  video.OwnerID,
			PlayUrl:  video.PlayURL,
			CoverUrl: video.CoverURL,
		},
	}, nil
}
