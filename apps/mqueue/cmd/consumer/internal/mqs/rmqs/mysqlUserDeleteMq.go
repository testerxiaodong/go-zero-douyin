package rmqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/svc"
	searchPb "go-zero-douyin/apps/search/cmd/rpc/pb"
	messageType "go-zero-douyin/common/message"
)

// MysqlUserDeleteMq 用户关注缓存
type MysqlUserDeleteMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMysqlUserDeleteMq(ctx context.Context, svcCtx *svc.ServiceContext) *MysqlUserDeleteMq {
	return &MysqlUserDeleteMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (v *MysqlUserDeleteMq) Consume(message string) error {
	// 获取消息内容：视频id
	var mysqlUserDeleteMessage messageType.MysqlUserDeleteMessage
	if err := json.Unmarshal([]byte(message), &mysqlUserDeleteMessage); err != nil {
		logx.WithContext(v.ctx).Error("mysqlUserDeleteMessage->Consume Unmarshal err : %v , val : %s", err, message)
		return err
	}
	// 删除视频文档
	_, err := v.svcCtx.SearchRpc.DeleteUser(v.ctx, &searchPb.DeleteUserDocumentReq{Id: mysqlUserDeleteMessage.UserId})
	if err != nil {
		return err
	}
	return nil
}
