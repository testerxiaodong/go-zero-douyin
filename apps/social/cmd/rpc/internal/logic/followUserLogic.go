package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
	"go-zero-douyin/common/message"
	"go-zero-douyin/common/xerr"
	"gorm.io/gorm"

	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowUserLogic {
	return &FollowUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FollowUser 关注功能
func (l *FollowUserLogic) FollowUser(in *pb.FollowUserReq) (*pb.FollowUserResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "follow user with empty param")
	}
	if in.GetUserId() == 0 || in.GetFollowerId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "follow user with empty follower_id or user_id")
	}

	// 用户不能关注自己
	if in.GetUserId() == in.GetFollowerId() {
		return nil, errors.Wrapf(xerr.NewErrMsg("不能对自己操作"), "req: %v", in)
	}

	// 查询数据库
	followQuery := l.svcCtx.Query.Follow
	follow, err := followQuery.WithContext(l.ctx).Where(followQuery.FollowerID.Eq(in.GetFollowerId())).Where(followQuery.UserID.Eq(in.GetUserId())).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "search user is alreaddy follow user from db failed, err: %v, follower_id: %d user_id: %d", err, in.GetFollowerId(), in.GetUserId())
	}

	// 已关注，直接返回成功，不做任何处理
	if follow != nil {
		return &pb.FollowUserResp{}, nil
	}

	// 未关注，插入关注记录
	newFollow := &model.Follow{}
	newFollow.FollowerID = in.GetFollowerId()
	newFollow.UserID = in.GetUserId()
	err = followQuery.WithContext(l.ctx).Create(newFollow)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "insert follow record failed, err: %v follower_id: %d user_id: %d", err, in.GetFollowerId(), in.GetUserId())
	}

	// 发布更新缓存消息
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
	return &pb.FollowUserResp{}, nil
}
