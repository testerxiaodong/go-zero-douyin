package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoUnlikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoUnlikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoUnlikeLogic {
	return &VideoUnlikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VideoUnlikeLogic) VideoUnlike(in *pb.VideoUnlikeReq) (*pb.VideoUnlikeResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unlike video with empty param")
	}
	if in.GetVideoId() == 0 || in.GetUserId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unlike video with empty video_id or user_id")
	}
	// 查询数据库
	like, err := l.svcCtx.LikeModel.FindOneByVideoIdUserId(l.ctx, in.GetVideoId(), in.GetUserId())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR),
			"find video is liked by user failed, err: %v", err)
	}
	// 没有记录，直接返回
	if like == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("用户尚未点赞"), "video_id: %d", in.GetVideoId())
	}
	// 有记录，但点赞状态为已取消，直接返回
	if like != nil && like.Status == xconst.LikeStateNo {
		return &pb.VideoUnlikeResp{}, nil
	}
	// 有记录且点赞状态为已点赞，更改状态，更改视频点赞数
	if like != nil && like.Status == xconst.LikeStateYes {
		if err := l.svcCtx.LikeModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {
			// 更新点赞状态
			like.Status = xconst.DelStateNo
			err := l.svcCtx.LikeModel.UpdateWithVersion(l.ctx, session, like)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR), "更新视频点赞状态失败")
			}
			// 更新点赞数
			likeCount, err := l.svcCtx.LikeCountModel.FindOneByVideoId(l.ctx, in.GetVideoId())
			if err != nil && !errors.Is(err, model.ErrNotFound) {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
					"查询视频点赞数失败, err: %v, video_id: %d", err, in.GetVideoId())
			}
			likeCount.LikeCount -= 1
			err = l.svcCtx.LikeCountModel.UpdateWithVersion(l.ctx, session, likeCount)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR),
					"更新视频点赞数失败, err: %v, video_id: %d", err, in.GetVideoId())
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return &pb.VideoUnlikeResp{}, nil
}
