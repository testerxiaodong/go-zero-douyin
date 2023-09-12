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

type AddTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTagLogic {
	return &AddTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddTag 新增标签
func (l *AddTagLogic) AddTag(in *pb.AddTagReq) (*pb.AddTagResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "add tag with empty param")
	}

	if len(in.GetName()) == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "add tag with empty name")
	}

	// 查询标签名是否已存在
	tag, err := l.svcCtx.TagDo.GetTagByName(l.ctx, in.GetName())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "数据库根据名称查询标签失败, err: %v, name: %s", err, in.GetName())
	}

	if tag != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("标签名已存在"), "新增标签失败， req: %v", in)
	}

	// 插入数据
	newTag := &model.Tag{}
	newTag.Name = in.GetName()

	err = l.svcCtx.TagDo.InsertTag(l.ctx, newTag)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "数据库新增标签失败, err: %v, name: %s", err, in.GetName())
	}

	// 删除redis缓存
	_, err = l.svcCtx.Redis.Delete(l.ctx, xconst.RedisVideoSection)
	// 缓存删除失败，发送kafka消息
	if err != nil {
		l.Logger.Errorf("redis标签信息删除失败, err: %v", err)
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
	return &pb.AddTagResp{}, nil
}
