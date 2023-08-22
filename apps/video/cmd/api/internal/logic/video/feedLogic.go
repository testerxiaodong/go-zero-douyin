package video

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	pb2 "go-zero-douyin/apps/comment/cmd/rpc/pb"
	pb3 "go-zero-douyin/apps/like/cmd/rpc/pb"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"sync"

	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var wg sync.WaitGroup

type FeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeedLogic) Feed(req *types.VideoFeedReq) (resp *types.VideoFeedResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := utils.GetValidator().ValidateZh(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}

	// 调用videorpc
	feedResp, err := l.svcCtx.VideoRpc.VideoFeed(l.ctx, &pb.VideoFeedReq{LastTimeStamp: utils.FromInt64TimeStampToProtobufTimeStamp(req.LastTimeStamp)})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}

	// 调用commentrpc获取视频评论数，调用likerpc获取视频点赞数
	if len(feedResp.Videos) > 0 {
		resp = &types.VideoFeedResp{Videos: make([]*types.Video, 0)}
		err = copier.Copy(resp, feedResp)
		if err != nil {
			return nil, errors.Wrapf(err, "copier feed resp failed: %v", feedResp)
		}
		wg.Add(len(feedResp.Videos))
		for i, video := range feedResp.Videos {
			index, currentVideo := i, video
			go func(index int, video *pb.Video) {
				defer wg.Done()
				commentCountResp, err := l.svcCtx.CommentRpc.GetCommentCountByVideoId(l.ctx, &pb2.GetCommentCountByVideoIdReq{VideoId: video.Id})
				if err != nil {
					logx.WithContext(l.ctx).Errorf("get video comment count by comment rpc failed, err: %v", err)
				}
				likeCountResp, err := l.svcCtx.LikeRpc.GetVideoLikeCountByVideoId(l.ctx, &pb3.GetVideoLikeCountByVideoIdReq{VideoId: video.Id})
				if err != nil {
					logx.WithContext(l.ctx).Errorf("get video like count by like rpc failed, err: %v", err)
				}
				resp.Videos[index].CommentCount = commentCountResp.Count
				resp.Videos[index].LikeCount = likeCountResp.LikeCount
			}(index, currentVideo)
		}
		wg.Wait()
	}

	return resp, nil
}
