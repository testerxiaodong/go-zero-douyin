package rmqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/svc"
	messageType "go-zero-douyin/common/message"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
)

// UserFollowedByUserMq 用户关注缓存
type UserFollowedByUserMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFollowedByUserMq(ctx context.Context, svcCtx *svc.ServiceContext) *UserFollowedByUserMq {
	return &UserFollowedByUserMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (v *UserFollowedByUserMq) Consume(message string) error {
	// 获取消息内容：视频id
	var userFollowedByUserMessage messageType.UserFollowedByUserMessage
	if err := json.Unmarshal([]byte(message), &userFollowedByUserMessage); err != nil {
		logx.WithContext(v.ctx).Error("userFollowUserMessage->Consume Unmarshal err : %v , val : %s", err, message)
		return err
	}
	logx.WithContext(v.ctx).Infof("获取到用户被关注消息，用户id: %d", userFollowedByUserMessage.UserId)
	// 删除视频缓存
	_, err := v.svcCtx.Redis.Del(utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowedByUserPrefix, userFollowedByUserMessage.UserId))
	if err != nil {
		// 少于重试最高次数，重新入队
		return err
	}
	logx.WithContext(v.ctx).Info("用户被关注的user_id集合缓存被删除")
	return nil
}
