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

type DelSectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelSectionLogic {
	return &DelSectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelSectionLogic) DelSection(in *pb.DelSectionReq) (*pb.DelSectionResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del section with empty param")
	}

	if in.GetId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del section with empty id")
	}

	// 删除数据
	_, err := l.svcCtx.SectionDo.DeleteSection(l.ctx, &model.Section{ID: in.GetId()})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "删除分区失败, err: %v, id: %d", err, in.GetId())
	}

	// 删除redis缓存
	_, err = l.svcCtx.Redis.Delete(l.ctx, xconst.RedisVideoSection)
	// 缓存删除失败，发送kafka消息
	if err != nil {
		body, err := json.Marshal(message.VideoSectionMessage{})
		if err != nil {
			panic(err)
		}
		err = l.svcCtx.Rabbit.Send("", "VideoSection", body)
		if err != nil {
			l.Logger.Errorf("发布删除视频分区缓存消息失败, err: %v", err)
			return nil, errors.Wrapf(xerr.NewErrMsg("发布删除分区缓存信息失败"), "err: %v", err)
		}
	}

	return &pb.DelSectionResp{}, nil
}
