package rmqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/svc"
	searchPb "go-zero-douyin/apps/search/cmd/rpc/pb"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	messageType "go-zero-douyin/common/message"
)

// MysqlUserUpdateMq 用户关注缓存
type MysqlUserUpdateMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMysqlUserUpdateMq(ctx context.Context, svcCtx *svc.ServiceContext) *MysqlUserUpdateMq {
	return &MysqlUserUpdateMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (v *MysqlUserUpdateMq) Consume(message string) error {
	// 获取消息内容：视频id
	var mysqlUserUpdateMessage messageType.MysqlUserUpdateMessage
	if err := json.Unmarshal([]byte(message), &mysqlUserUpdateMessage); err != nil {
		logx.WithContext(v.ctx).Error("mysqlUserUpdateMessage->Consume Unmarshal err : %v , val : %s", err, message)
		return err
	}
	userInfoResp, err := v.svcCtx.UserRpc.GetUserInfo(v.ctx, &pb.GetUserInfoReq{Id: mysqlUserUpdateMessage.UserId})
	if err != nil {
		return err
	}
	followerCountResp, err := v.svcCtx.SocialRpc.GetUserFollowerCount(v.ctx, &socialPb.GetUserFollowerCountReq{UserId: mysqlUserUpdateMessage.UserId})
	if err != nil {
		return err
	}
	followCountResp, err := v.svcCtx.SocialRpc.GetUserFollowCount(v.ctx, &socialPb.GetUserFollowCountReq{UserId: mysqlUserUpdateMessage.UserId})
	if err != nil {
		return err
	}
	_, err = v.svcCtx.SearchRpc.SyncUserInfo(v.ctx, &searchPb.SyncUserInfoReq{User: &searchPb.User{
		Id:            userInfoResp.User.Id,
		Username:      userInfoResp.User.Username,
		FollowerCount: followerCountResp.FollowerCount,
		FollowCount:   followCountResp.FollowCount,
	}})
	if err != nil {
		return err
	}
	return nil
}
