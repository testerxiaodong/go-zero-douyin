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

type VideoLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoLikeLogic {
	return &VideoLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VideoLikeLogic) VideoLike(in *pb.VideoLikeReq) (*pb.VideoLikeResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "like video with empty param")
	}
	if in.GetUserId() == 0 || in.GetVideoId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR),
			"like video with empty user_id or video_id")
	}
	// 查询记录是否存在
	like, err := l.svcCtx.LikeModel.FindOneByVideoIdUserIdIsDelete(l.ctx, in.GetVideoId(), in.GetUserId(), xconst.DelStateNo)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR),
			"find video is already liked by user failed, err: %v", err)
	}
	// 记录存在，且状态为已点赞，直接返回
	if like != nil && like.Status == xconst.LikeStateYes {
		return &pb.VideoLikeResp{}, nil
	}
	// 记录存在，且状态为未点赞，更新状态，更新视频点赞数
	if like != nil && like.Status == xconst.LikeStateNo {
		if err := l.svcCtx.LikeModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {
			// 更新点赞状态
			like.Status = xconst.LikeStateYes
			err := l.svcCtx.LikeModel.UpdateWithVersion(l.ctx, session, like)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR), "更新")
			}
			// 更新视频点赞数
			likeCount, err := l.svcCtx.LikeCountModel.FindOneByVideoIdIsDelete(l.ctx, in.GetVideoId(), xconst.DelStateNo)
			if err != nil && !errors.Is(err, model.ErrNotFound) {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
					"查询视频点赞数失败, err: %v, video_id: %d", err, in.GetVideoId())
			}
			likeCount.LikeCount += 1
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
	// 记录不存在，插入记录，更新点赞数
	if like == nil {
		if err := l.svcCtx.LikeModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {
			// 插入记录
			newLike := &model.Like{}
			newLike.UserId = in.GetUserId()
			newLike.VideoId = in.GetVideoId()
			newLike.Status = xconst.LikeStateYes
			_, err := l.svcCtx.LikeModel.Insert(l.ctx, session, newLike)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR),
					"插入点赞记录失败, err: %v, user_id: %d, video_id: %d", err, in.GetUserId(), in.GetVideoId())
			}
			// 更新点赞数
			likeCount, err := l.svcCtx.LikeCountModel.FindOneByVideoIdIsDelete(l.ctx, in.GetVideoId(), xconst.DelStateNo)
			if err != nil && !errors.Is(err, model.ErrNotFound) {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
					"查询视频点赞数失败, err: %v, video_id: %d", err, in.GetVideoId())
			}
			// 记录不存在，插入一条记录
			if likeCount == nil {
				newLikeCount := &model.LikeCount{}
				newLikeCount.LikeCount = 1
				newLikeCount.VideoId = in.GetVideoId()
				_, err = l.svcCtx.LikeCountModel.Insert(l.ctx, session, newLikeCount)
				if err != nil {
					return errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR),
						"插入视频点赞数记录失败, err: %v, video_id: %d", err, in.GetVideoId())
				}
			} else {
				// 记录存在，更新点赞数+1
				likeCount.LikeCount += 1
				err := l.svcCtx.LikeCountModel.UpdateWithVersion(l.ctx, session, likeCount)
				if err != nil {
					return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR),
						"更新视频点赞数失败, err; %v, video_id: %d", err, in.GetVideoId())
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return &pb.VideoLikeResp{}, nil
}
