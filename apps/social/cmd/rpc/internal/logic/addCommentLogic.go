package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
	"go-zero-douyin/common/message"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCommentLogic {
	return &AddCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddComment 评论相关功能
func (l *AddCommentLogic) AddComment(in *pb.AddCommentReq) (*pb.AddCommentResp, error) {
	// todo: add your logic here and delete this line
	// 校验参数
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Add comment with empty param")
	}
	if in.GetVideoId() == 0 || in.GetUserId() == 0 || len(in.GetContent()) == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Add comment with empty video_id or user_id or content")
	}

	// 复制数据
	comment := &model.Comment{}
	comment.UserID = in.GetUserId()
	comment.VideoID = in.GetVideoId()
	comment.Content = in.GetContent()

	// 插入评论
	err := l.svcCtx.CommentDo.InsertComment(l.ctx, comment)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "insert comment failed: %v", err)
	}

	// 删除缓存
	if _, err := l.svcCtx.Redis.Delete(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoCommentPrefix, in.GetVideoId())); err != nil {
		// 删除缓存失败，发布消息异步处理
		body, err := json.Marshal(message.VideoCommentMessage{VideoId: in.GetVideoId()})
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "marshal video comment count message failed: %v", err)
		}
		err = l.svcCtx.Rabbit.Send("", "VideoCommentMq", body)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "publish video comment count message failed: %v", err)
		}
	}
	return &pb.AddCommentResp{}, nil
}
