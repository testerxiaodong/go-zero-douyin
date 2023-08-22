package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/like/cmd/rpc/internal/model"
	"go-zero-douyin/common/message"
	"go-zero-douyin/common/xerr"
	"gorm.io/gorm"

	"go-zero-douyin/apps/like/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/like/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoLikeLogic {
	return &VideoLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VideoLikeLogic) VideoLike(in *pb.VideoLikeReq) (*pb.VideoLikeResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "like video with empty param")
	}
	if in.GetUserId() == 0 || in.GetVideoId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "like video with empty user_id or video_id")
	}

	// 查询是否已点赞
	likeQuery := l.svcCtx.Query.Like
	like, err := likeQuery.WithContext(l.ctx).Where(likeQuery.VideoID.Eq(in.GetVideoId())).Where(likeQuery.UserID.Eq(in.GetUserId())).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "find video is already liked by user failed, err: %v", err)
	}
	if like != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("video already liked by user"), "video_id: %d, user_id: %d", in.GetVideoId(), in.GetUserId())
	}

	// 插入数据库
	newLike := &model.Like{}
	newLike.VideoID = in.GetVideoId()
	newLike.UserID = in.GetUserId()
	err = likeQuery.WithContext(l.ctx).Create(newLike)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_INSERT_ERR), "insert video like failed, err: %v", err)
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
