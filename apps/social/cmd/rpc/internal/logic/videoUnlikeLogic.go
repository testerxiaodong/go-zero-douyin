package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-zero-douyin/common/message"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
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

func (l *VideoUnlikeLogic) VideoUnlike(in *pb.VideoUnlikeReq) (*pb.VideoUnlikeResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unlike video with empty param")
	}
	if in.GetVideoId() == 0 || in.GetUserId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unlike video with empty video_id or user_id")
	}

	// 查询数据库
	like, err := l.svcCtx.LikeDo.GetLikeByVideoIdAndUserId(l.ctx, in.GetVideoId(), in.GetUserId())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "find video is liked by user failed, err: %v", err)
	}

	// 没有记录，直接返回
	if like == nil {
		return &pb.VideoUnlikeResp{}, nil
	}

	// 删除点赞记录
	_, err = l.svcCtx.LikeDo.DeleteLike(l.ctx, like)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_DELETE_ERR), "delete user like video record failed, err: %v", err)
	}

	// 删除用户点赞视频id集合缓存
	if _, err := l.svcCtx.Redis.Delete(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserLikeVideoPrefix, in.GetUserId())); err != nil {
		// 删除缓存失败，发布消息异步处理
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
		// 删除缓存失败，发布消息异步处理
		videoUserBody, err := json.Marshal(message.VideoLikedByUserMessage{VideoId: in.GetVideoId()})
		if err != nil {
			panic(err)
		}

		err = l.svcCtx.Rabbit.Send("", "VideoLikedByUserMq", videoUserBody)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("publish video liked by user message failed"), "user_id: %d", in.GetUserId())
		}
	}

	return &pb.VideoUnlikeResp{}, nil
}
