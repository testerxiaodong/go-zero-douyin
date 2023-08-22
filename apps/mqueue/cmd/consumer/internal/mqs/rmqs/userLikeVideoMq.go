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

// UserLikeVideoMq 视频评论缓存信息
type UserLikeVideoMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLikeVideoMq(ctx context.Context, svcCtx *svc.ServiceContext) *UserLikeVideoMq {
	return &UserLikeVideoMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (v *UserLikeVideoMq) Consume(message string) error {
	// 获取消息内容：视频id
	var userLikeVideoMessage messageType.UserLikeVideoMessage
	if err := json.Unmarshal([]byte(message), &userLikeVideoMessage); err != nil {
		logx.WithContext(v.ctx).Error("userLikeVideoMessage->Consume Unmarshal err : %v , val : %s", err, message)
		return err
	}
	// 删除视频缓存
	_, err := v.svcCtx.Redis.Del(utils.GetRedisKeyWithPrefix(xconst.RedisUserLikeVideoPrefix, userLikeVideoMessage.UserId))
	if err != nil {
		// 少于重试最高次数，重新入队
		return err
	}
	return nil
}
