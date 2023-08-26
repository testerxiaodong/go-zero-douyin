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
	followQuery := l.svcCtx.Query.Follow
	follow, err := followQuery.WithContext(l.ctx).Where(followQuery.FollowerID.Eq(in.GetFollowerId())).Where(followQuery.UserID.Eq(in.GetUserId())).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "find follow record by follower_id and user_id failed, err: %v, follower_id: %d, user_id: %d", err, in.GetFollowerId(), in.GetUserId())
	}

	// 没有记录，直接返回删除成功
	if follow == nil {
		return &pb.UnfollowUserResp{}, nil
	}

	// 删除数据库记录，取消关注
	_, err = followQuery.WithContext(l.ctx).Delete(follow)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "delete follow record failed, err: %v, follower_id: %d, user_id: %d", err, in.GetFollowerId(), in.GetUserId())
	}

	// 发布删除缓存信息
	userFollowUserMessage := message.UserFollowUserMessage{FollowerId: in.GetFollowerId()}
	userFollowedByUserMessage := message.UserFollowedByUserMessage{UserId: in.GetUserId()}

	userFollowUserBody, err := json.Marshal(userFollowUserMessage)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("序列化userFollowUserMessage失败"), "err: %v data: %v", err, userFollowUserMessage)
	}
	userFollowedByUserBody, err := json.Marshal(userFollowedByUserMessage)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("序列化userFollowedByUserMessage失败"), "err: %v data: %v", err, userFollowedByUserMessage)
	}
	err = l.svcCtx.Rabbit.Send("", "UserFollowUserMq", userFollowUserBody)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("发布userFollowUserMessage失败"), "err: %v", err)
	}
	err = l.svcCtx.Rabbit.Send("", "UserFollowedByUserMq", userFollowedByUserBody)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("发布userFollowUserMessage失败"), "err: %v", err)
	}
	return &pb.UnfollowUserResp{}, nil
}
