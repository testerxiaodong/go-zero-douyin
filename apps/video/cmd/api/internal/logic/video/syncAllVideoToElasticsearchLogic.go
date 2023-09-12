package video

import (
	"context"
	"github.com/pkg/errors"
	pbSocial "go-zero-douyin/apps/social/cmd/rpc/pb"
	pbVideo "go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-douyin/apps/video/cmd/api/internal/svc"
)

type SyncAllVideoToElasticsearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncAllVideoToElasticsearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncAllVideoToElasticsearchLogic {
	return &SyncAllVideoToElasticsearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncAllVideoToElasticsearchLogic) SyncAllVideoToElasticsearch() error {
	// todo: add your logic here and delete this line
	// 获取所有视频信息
	videoList, err := l.svcCtx.VideoRpc.GetAllVideo(l.ctx, &pbVideo.GetAllVideoReq{})
	if err != nil {
		return errors.Wrap(err, "获取所有视频信息失败")
	}

	// 调用socialRpc获取视频点赞数和评论数
	if len(videoList.Videos) > 0 {
		wg.Add(len(videoList.Videos))
		for i, video := range videoList.Videos {
			index, currentVideo := i, video
			go func(index int, video *pbVideo.VideoInfo) {
				videoInfo := &pbVideo.VideoDetailInfo{}
				videoInfo.Id = video.Id
				videoInfo.Title = video.Title
				videoInfo.SectionId = video.SectionId
				videoInfo.Tags = video.Tags
				videoInfo.OwnerId = video.OwnerId
				videoInfo.PlayUrl = video.PlayUrl
				videoInfo.CoverUrl = video.CoverUrl
				videoInfo.CreateTime = video.CreateTime
				videoInfo.UpdateTime = video.UpdateTime
				defer wg.Done()
				countResp, err := l.svcCtx.SocialRpc.GetCommentCountByVideoId(l.ctx, &pbSocial.GetCommentCountByVideoIdReq{VideoId: video.Id})
				if err != nil {
					logx.WithContext(l.ctx).Errorf("get video comment count by comment rpc failed, err: %v", err)
				} else {
					videoInfo.CommentCount = countResp.Count
				}
				likeCountResp, err := l.svcCtx.SocialRpc.GetVideoLikedCountByVideoId(l.ctx, &pbSocial.GetVideoLikedCountByVideoIdReq{VideoId: video.Id})
				if err != nil {
					logx.WithContext(l.ctx).Errorf("get video like count by like rpc failed, err: %v", err)
				} else {
					videoInfo.LikeCount = likeCountResp.LikeCount
				}
				_, err = l.svcCtx.VideoRpc.SyncVideoInfoToElasticsearch(l.ctx, &pbVideo.SyncVideoInfoToElasticsearchReq{VideoInfo: videoInfo})
				if err != nil {
					l.Logger.Errorf("VideoRpc.SyncVideoInfoToElasticsearch失败, err: %v", err)
				}
			}(index, currentVideo)
		}
		wg.Wait()
	}

	return nil
}
