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

type UnfollowUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnfollowUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfollowUserLogic {
	return &UnfollowUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnfollowUserLogic) UnfollowUser(in *pb.UnfollowUserReq) (*pb.UnfollowUserResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unfollow user with empty param")
	}
	if in.GetUserId() == 0 || in.GetFollowerId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unfollow user with empty follower_id or user_id")
	}

	// 不能对自己操作
	if in.GetUserId() == in.GetFollowerId() {
		return nil, errors.Wrapf(xerr.NewErrMsg("不能对自己操作"), "req: %v", in)
	}

	// 查询数据库
	follow, err := l.svcCtx.FollowDo.GetFollowByFollowerIdAndUserId(l.ctx, in.GetFollowerId(), in.GetUserId())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "find follow record by follower_id and user_id failed, err: %v, follower_id: %d, user_id: %d", err, in.GetFollowerId(), in.GetUserId())
	}

	// 没有记录，直接返回删除成功
	if follow == nil {
		return &pb.UnfollowUserResp{}, nil
	}

	// 删除数据库记录，取消关注
	_, err = l.svcCtx.FollowDo.DeleteFollow(l.ctx, follow)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "delete follow record failed, err: %v, follower_id: %d, user_id: %d", err, in.GetFollowerId(), in.GetUserId())
	}

	// 删除用户关注id集合缓存
	if _, err := l.svcCtx.Redis.Delete(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowUserPrefix, in.GetFollowerId())); err != nil {
		// 删除失败，发布异步处理消息
		userFollowUserMessage := message.UserFollowUserMessage{FollowerId: in.GetFollowerId()}

		userFollowUserBody, err := json.Marshal(userFollowUserMessage)
		if err != nil {
			panic(err)
		}

		err = l.svcCtx.Rabbit.Send("", "UserFollowUserMq", userFollowUserBody)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("发布userFollowUserMessage失败"), "err: %v", err)
		}
	}

	// 删除用户粉丝id集合缓存
	if _, err := l.svcCtx.Redis.Delete(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowedByUserPrefix, in.GetUserId())); err != nil {
		// 删除失败，发布异步处理消息
		userFollowedByUserMessage := message.UserFollowedByUserMessage{UserId: in.GetUserId()}

		userFollowedByUserBody, err := json.Marshal(userFollowedByUserMessage)
		if err != nil {
			panic(err)
		}

		err = l.svcCtx.Rabbit.Send("", "UserFollowedByUserMq", userFollowedByUserBody)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("发布userFollowedByUserMessage失败"), "err: %v", err)
		}
	}
	// 发布更新es用户的消息
	userMsg, _ := json.Marshal(message.MysqlUserUpdateMessage{UserId: in.GetUserId()})
	err = l.svcCtx.Rabbit.Send("", "MysqlUserUpdateMq", userMsg)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_UPDATE_ERR), "req: %v, err: %v", in, err)
	}
	followerMsg, _ := json.Marshal(message.MysqlUserUpdateMessage{UserId: in.GetFollowerId()})
	err = l.svcCtx.Rabbit.Send("", "MysqlUserUpdateMq", followerMsg)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_UPDATE_ERR), "req: %v, err: %v", in, err)
	}

	return &pb.UnfollowUserResp{}, nil
}
