package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
	"go-zero-douyin/common/message"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"
	"gorm.io/gorm"

	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"

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
	like, err := l.svcCtx.LikeDo.GetLikeByVideoIdAndUserId(l.ctx, in.GetVideoId(), in.GetUserId())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "find video is already liked by user failed, err: %v", err)
	}

	// 已经点赞，直接返回
	if like != nil {
		return &pb.VideoLikeResp{}, nil
	}

	// 插入数据库
	newLike := &model.Like{}
	newLike.VideoID = in.GetVideoId()
	newLike.UserID = in.GetUserId()
	err = l.svcCtx.LikeDo.InsertLike(l.ctx, newLike)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_INSERT_ERR), "insert video like failed, err: %v", err)
	}

	// 删除用户点赞视频id集合缓存
	if _, err := l.svcCtx.Redis.Delete(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserLikeVideoPrefix, in.GetUserId())); err != nil {
		userVideoBody, err := json.Marshal(message.UserLikeVideoMessage{UserId: in.GetUserId()})
		if err != nil {
			panic(err)
		}
		err = l.svcCtx.Rabbit.Send("", "UserLikeVideoMq", userVideoBody)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("publish user like video message failed"), "video_id: %d", in.GetVideoId())
		}
	}

	// 删除视频被点赞用户id集合缓存
	if _, err := l.svcCtx.Redis.Delete(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoLikedByUserPrefix, in.GetVideoId())); err != nil {
		videoUserBody, err := json.Marshal(message.VideoLikedByUserMessage{VideoId: in.GetVideoId()})
		if err != nil {
			panic(err)
		}

		err = l.svcCtx.Rabbit.Send("", "VideoLikedByUserMq", videoUserBody)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("publish video liked by user message failed"), "user_id: %d", in.GetUserId())
		}
	}

	return &pb.VideoLikeResp{}, nil
}
