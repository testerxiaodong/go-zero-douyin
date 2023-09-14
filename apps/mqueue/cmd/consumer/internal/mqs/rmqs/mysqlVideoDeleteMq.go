package rmqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/svc"
	searchPb "go-zero-douyin/apps/search/cmd/rpc/pb"
	messageType "go-zero-douyin/common/message"
)

// MysqlVideoDeleteMq 用户关注缓存
type MysqlVideoDeleteMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMysqlVideoDeleteMq(ctx context.Context, svcCtx *svc.ServiceContext) *MysqlVideoDeleteMq {
	return &MysqlVideoDeleteMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (v *MysqlVideoDeleteMq) Consume(message string) error {
	// 获取消息内容：视频id
	var mysqlVideoDeleteMessage messageType.MysqlVideoDeleteMessage
	if err := json.Unmarshal([]byte(message), &mysqlVideoDeleteMessage); err != nil {
		logx.WithContext(v.ctx).Error("mysqlVideoDeleteMessage->Consume Unmarshal err : %v , val : %s", err, message)
		return err
	}
	// 删除视频文档
	_, err := v.svcCtx.SearchRpc.DeleteVideo(v.ctx, &searchPb.DeleteVideoDocumentReq{Id: mysqlVideoDeleteMessage.VideoId})
	if err != nil {
		return err
	}
	return nil
}
