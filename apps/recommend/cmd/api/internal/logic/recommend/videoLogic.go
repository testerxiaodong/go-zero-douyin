package recommend

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/recommend/cmd/rpc/pb"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
	videoPb "go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/ctxdata"
	"go-zero-douyin/common/xerr"
	"golang.org/x/sync/errgroup"

	"go-zero-douyin/apps/recommend/cmd/api/internal/svc"
	"go-zero-douyin/apps/recommend/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoLogic {
	return &VideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoLogic) Video(req *types.VideoRecommendReq) (resp *types.VideoRecommendResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}
	// 获取用户id
	uid := ctxdata.GetUidFromCtx(l.ctx)
	// 调用RecommendRpc
	videoRecommendSectionResp, err := l.svcCtx.RecommendRpc.VideoRecommendSection(l.ctx,
		&pb.VideoRecommendSectionReq{UserId: uid, SectionId: req.SectionId, Count: req.Count})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}
	// 获取视频信息
	resp = &types.VideoRecommendResp{Videos: make([]*types.Video, 0)}
	if len(videoRecommendSectionResp.GetVideoIds()) > 0 {
		for _, videoId := range videoRecommendSectionResp.VideoIds {
			G, _ := errgroup.WithContext(l.ctx)
			video := &types.Video{}
			// 调用VideoRpc获取视频基本信息
			G.Go(func() error {
				videoInfoResp, err := l.svcCtx.VideoRpc.GetVideoById(l.ctx, &videoPb.GetVideoByIdReq{Id: videoId})
				if err != nil {
					return err
				}
				err = copier.Copy(&video, videoInfoResp.Video)
				if err != nil {
					return err
				}
				return nil
			})
			// 调用SocialRpc获取视频点赞数
			G.Go(func() error {
				likedCountByVideoIdResp, err := l.svcCtx.SocialRpc.GetVideoLikedCountByVideoId(l.ctx,
					&socialPb.GetVideoLikedCountByVideoIdReq{VideoId: videoId})
				if err != nil {
					return err
				}
				video.LikeCount = likedCountByVideoIdResp.LikeCount
				return nil
			})
			// 调用SocialRpc获取视频评论数
			G.Go(func() error {
				likedCountByVideoIdResp, err := l.svcCtx.SocialRpc.GetCommentCountByVideoId(l.ctx,
					&socialPb.GetCommentCountByVideoIdReq{VideoId: videoId})
				if err != nil {
					return err
				}
				video.CommentCount = likedCountByVideoIdResp.Count
				return nil
			})
			err = G.Wait()
			if err != nil {
				return nil, errors.Wrapf(err, "req: %v", req)
			}
			resp.Videos = append(resp.Videos, video)
		}
	}
	return resp, nil
}
