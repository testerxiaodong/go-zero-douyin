package video

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	pbSocial "go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.UserVideoListReq) (resp *types.UserVideoListResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}

	// 调用videoRpc
	videoList, err := l.svcCtx.VideoRpc.UserVideoList(l.ctx, &pb.UserVideoListReq{UserId: req.UserId})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}

	// 调用socialRpc
	if len(videoList.Videos) > 0 {
		resp = &types.UserVideoListResp{Videos: make([]*types.VideoInfo, 0)}
		err = copier.Copy(resp, videoList)
		if err != nil {
			return nil, errors.Wrapf(err, "copier feed resp failed: %v", videoList)
		}
		wg.Add(len(videoList.Videos))
		for i, video := range videoList.Videos {
			index, currentVideo := i, video
			go func(index int, video *pb.VideoInfo) {
				defer wg.Done()
				countResp, err := l.svcCtx.SocialRpc.GetCommentCountByVideoId(l.ctx, &pbSocial.GetCommentCountByVideoIdReq{VideoId: video.Id})
				if err != nil {
					resp.Videos[index].CommentCount = 0
					logx.WithContext(l.ctx).Errorf("get video comment count by comment rpc failed, err: %v", err)
				} else {
					resp.Videos[index].CommentCount = countResp.Count
				}
				likeCountResp, err := l.svcCtx.SocialRpc.GetVideoLikedCountByVideoId(l.ctx, &pbSocial.GetVideoLikedCountByVideoIdReq{VideoId: video.Id})
				if err != nil {
					resp.Videos[index].LikeCount = 0
					logx.WithContext(l.ctx).Errorf("get video like count by like rpc failed, err: %v", err)
				} else {
					resp.Videos[index].LikeCount = likeCountResp.LikeCount
				}
			}(index, currentVideo)
		}
		wg.Wait()
	}
	return resp, nil
}
