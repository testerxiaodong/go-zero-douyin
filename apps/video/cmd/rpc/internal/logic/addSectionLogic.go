package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/common/message"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"
	"gorm.io/gorm"

	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddSectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSectionLogic {
	return &AddSectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddSection 添加分区
func (l *AddSectionLogic) AddSection(in *pb.AddSectionReq) (*pb.AddSectionResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "add section with empty param")
	}

	if len(in.GetName()) == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "add section with empty name")
	}

	// 查询分区名是否已存在
	section, err := l.svcCtx.SectionDo.GetSectionByName(l.ctx, in.GetName())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "数据库根据名称查询分区失败, err: %v, name: %s", err, in.GetName())
	}

	if section != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("分区名已存在"), "新增分区失败， req: %v", in)
	}

	// 插入数据
	newSection := &model.Section{}
	newSection.Name = in.GetName()

	err = l.svcCtx.SectionDo.InsertSection(l.ctx, newSection)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "数据库新增分区失败, err: %v, name: %s", err, in.GetName())
	}

	// 删除redis缓存
	_, err = l.svcCtx.Redis.Delete(l.ctx, xconst.RedisVideoSection)
	// 缓存删除失败，发送kafka消息
	if err != nil {
		l.Logger.Errorf("redis分区信息删除失败, err: %v", err)
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

	return &pb.AddSectionResp{}, nil
}
