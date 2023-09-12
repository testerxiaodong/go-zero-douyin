package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/common/message"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelTagLogic {
	return &DelTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelTagLogic) DelTag(in *pb.DelTagReq) (*pb.DelTagResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del tag with empty param")
	}

	if in.GetId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del tag with empty id")
	}

	// 删除数据
	_, err := l.svcCtx.TagDo.DeleteTag(l.ctx, &model.Tag{ID: in.GetId()})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "删除标签失败, err: %v, id: %d", err, in.GetId())
	}

	// 删除redis缓存
	_, err = l.svcCtx.Redis.Delete(l.ctx, xconst.RedisVideoSection)
	// 缓存删除失败，发送kafka消息
	if err != nil {
		body, err := json.Marshal(message.VideoTagMessage{})
		if err != nil {
			panic(err)
		}
		err = l.svcCtx.Rabbit.Send("", "VideoTag", body)
		if err != nil {
			l.Logger.Errorf("发布删除视频标签缓存消息失败, err: %v", err)
			return nil, errors.Wrapf(xerr.NewErrMsg("发布删除标签缓存信息失败"), "err: %v", err)
		}
	}

	return &pb.DelTagResp{}, nil
}
