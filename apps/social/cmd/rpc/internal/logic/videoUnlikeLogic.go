package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-zero-douyin/common/message"
	"go-zero-douyin/common/xerr"
	"gorm.io/gorm"

	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoUnlikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoUnlikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoUnlikeLogic {
	return &VideoUnlikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VideoUnlikeLogic) VideoUnlike(in *pb.VideoUnlikeReq) (*pb.VideoLikeResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unlike video with empty param")
	}
	if in.GetVideoId() == 0 || in.GetUserId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unlike video with empty video_id or user_id")
	}

	// 查询数据库
	likeQuery := l.svcCtx.Query.Like
	like, err := likeQuery.WithContext(l.ctx).Where(likeQuery.VideoID.Eq(in.GetVideoId())).Where(likeQuery.UserID.Eq(in.GetUserId())).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "find video is liked by user failed, err: %v", err)
	}
	if like == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("video is not liked by user"), "video_id: %d user_id: %v", in.GetVideoId(), in.GetUserId())
	}

	// 删除点赞记录
	_, err = likeQuery.WithContext(l.ctx).Delete(like)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_DELETE_ERR), "delete user like video record failed, err: %v", err)
	}

	// 发布消息，异步删除缓存
	userVideoBody, err := json.Marshal(message.UserLikeVideoMessage{UserId: in.GetUserId()})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "marshal user like video message failed, err: %v", err)
	}
	videoUserBody, err := json.Marshal(message.VideoLikedByUserMessage{VideoId: in.GetVideoId()})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "marshal video liked by user message failed, err: %v", err)
	}
	err = l.svcCtx.Rabbit.Send("", "UserLikeVideoMq", userVideoBody)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("publish user like video message failed"), "video_id: %d", in.GetVideoId())
	}
	err = l.svcCtx.Rabbit.Send("", "VideoLikedByUserMq", videoUserBody)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("publish video liked by user message failed"), "user_id: %d", in.GetUserId())
	}

	return &pb.VideoLikeResp{}, nil
}
