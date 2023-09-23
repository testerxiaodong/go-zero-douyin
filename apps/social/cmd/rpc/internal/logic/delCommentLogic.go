package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

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
	comment, err := l.svcCtx.CommentModel.FindOne(l.ctx, in.GetCommentId())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "find comment by id failed, err: %v", err)
	}

	// 评论不存在
	if comment == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("评论不存在"), "comment not found, id: %d", in.GetCommentId())
	}

	// 评论非该用户发布
	if comment.UserId != in.GetUserId() {
		return nil, errors.Wrapf(xerr.NewErrMsg("评论非该用户发布，无法删除"), "comment_id: %d", in.GetCommentId())
	}

	// 删除评论
	if err := l.svcCtx.CommentModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {
		// 删除评论
		err = l.svcCtx.CommentModel.DeleteSoft(l.ctx, nil, comment)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "del comment failed, err: %v", err)
		}
		// 更新评论数：-1
		commentCount, err := l.svcCtx.CommentCountModel.FindOneByVideoId(l.ctx, comment.VideoId)
		if err != nil && !errors.Is(err, model.ErrNotFound) {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"查询视频评论数记录失败, err: %v, video_id: %d", err, comment.VideoId)
		}
		commentCount.CommentCount -= 1
		err = l.svcCtx.CommentCountModel.UpdateWithVersion(l.ctx, session, commentCount)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR),
				"更新视频评论数失败, err: %v, video_id: %d", err, comment.VideoId)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &pb.DelCommentResp{}, nil
}
