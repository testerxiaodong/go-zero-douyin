package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-zero-douyin/common/message"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"
	"gorm.io/gorm"

	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelCommentLogic {
	return &DelCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelCommentLogic) DelComment(in *pb.DelCommentReq) (*pb.DelCommentResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del comment with empty param")
	}
	if in.GetUserId() == 0 || in.GetCommentId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del comment with empty user_id or comment_id")
	}

	// 查询数据库
	comment, err := l.svcCtx.CommentDo.GetCommentById(l.ctx, in.GetCommentId())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "find comment by id failed, err: %v", err)
	}

	// 评论不存在
	if comment == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "comment not found, id: %d", in.GetCommentId())
	}

	// 评论非该用户发布
	if comment.UserID != in.GetUserId() {
		return nil, errors.Wrapf(xerr.NewErrMsg("评论非该用户发布，无法删除"), "comment_id: %d", in.GetCommentId())
	}

	// 删除评论
	_, err = l.svcCtx.CommentDo.DeleteComment(l.ctx, comment)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "del comment failed, err: %v", err)
	}

	// 删除缓存
	if _, err := l.svcCtx.Redis.Delete(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoCommentPrefix, comment.VideoID)); err != nil {
		// 删除缓存失败，发布消息异步处理
		body, err := json.Marshal(message.VideoCommentMessage{VideoId: comment.VideoID})
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "marshal video comment count message failed: %v", err)
		}
		err = l.svcCtx.Rabbit.Send("", "VideoCommentMq", body)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "publish video comment count message failed: %v", err)
		}
	}

	// 发布更新es视频文档的消息
	msg, _ := json.Marshal(message.MysqlVideoUpdateMessage{VideoId: comment.VideoID})
	err = l.svcCtx.Rabbit.Send("", "MysqlVideoUpdateMq", msg)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_UPDATE_ERR), "req: %v, err: %v", in, err)
	}

	return &pb.DelCommentResp{}, nil
}
