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

// UserFollowUserMq 用户关注缓存
type UserFollowUserMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFollowUserMq(ctx context.Context, svcCtx *svc.ServiceContext) *UserFollowUserMq {
	return &UserFollowUserMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (v *UserFollowUserMq) Consume(message string) error {
	// 获取消息内容：视频id
	var userFollowUserMessage messageType.UserFollowUserMessage
	if err := json.Unmarshal([]byte(message), &userFollowUserMessage); err != nil {
		logx.WithContext(v.ctx).Error("userFollowUserMessage->Consume Unmarshal err : %v , val : %s", err, message)
		return err
	}
	// 删除视频缓存
	_, err := v.svcCtx.Redis.Del(utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowUserPrefix, userFollowUserMessage.FollowerId))
	if err != nil {
		// 少于重试最高次数，重新入队
		return err
	}
	return nil
}
