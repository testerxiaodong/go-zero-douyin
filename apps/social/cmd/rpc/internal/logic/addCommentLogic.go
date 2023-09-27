package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
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
	_ = copier.Copy(comment, in)
	// 插入评论，增加视频评论数
	if err := l.svcCtx.CommentModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {
		// 插入评论记录
		_, err := l.svcCtx.CommentModel.Insert(l.ctx, nil, comment)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "insert comment failed: %v", err)
		}
		commentCountRecord, err := l.svcCtx.CommentCountModel.FindOneByVideoIdIsDelete(l.ctx, in.GetVideoId(), xconst.DelStateNo)
		// 查询失败
		if err != nil && !errors.Is(err, model.ErrNotFound) {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"查询视频的评论数失败, err: %v, video_id: %d", err, in.GetVideoId())
		}
		// 记录不存在，插入一条评论数为1的记录
		if commentCountRecord == nil {
			commentCount := &model.CommentCount{}
			commentCount.VideoId = in.GetVideoId()
			commentCount.CommentCount = 1
			_, err = l.svcCtx.CommentCountModel.Insert(l.ctx, session, commentCount)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "插入视频评论数记录失败")
			}
		}
		// 记录存在，更新记录的评论数：+1
		if commentCountRecord != nil {
			commentCountRecord.CommentCount += 1
			err := l.svcCtx.CommentCountModel.UpdateWithVersion(l.ctx, session, commentCountRecord)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR), "更新视频评论数失败, err: %v, video_id: %d", err, in.GetVideoId())
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &pb.AddCommentResp{}, nil
}
