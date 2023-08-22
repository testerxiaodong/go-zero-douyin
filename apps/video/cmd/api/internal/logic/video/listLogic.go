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
	if validateResult := utils.GetValidator().ValidateZh(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}
	// 调用videorpc
	videoList, err := l.svcCtx.VideoRpc.UserVideoList(l.ctx, &pb.UserVideoListReq{UserId: req.UserId})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}
	// 调用commentrpc
	if len(videoList.Videos) > 0 {
		resp = &types.UserVideoListResp{Videos: make([]*types.Video, 0)}
		err = copier.Copy(resp, videoList)
		if err != nil {
			return nil, errors.Wrapf(err, "copier feed resp failed: %v", videoList)
		}
		wg.Add(len(videoList.Videos))
		for i, video := range videoList.Videos {
			index, currentVideo := i, video
			go func(index int, video *pb.Video) {
				defer wg.Done()
				countResp, err := l.svcCtx.CommentRpc.GetCommentCountByVideoId(l.ctx, &pb2.GetCommentCountByVideoIdReq{VideoId: video.Id})
				if err != nil {
					logx.WithContext(l.ctx).Errorf("get video comment count by comment rpc failed, err: %v", err)
				}
				likeCountResp, err := l.svcCtx.LikeRpc.GetVideoLikeCountByVideoId(l.ctx, &pb3.GetVideoLikeCountByVideoIdReq{VideoId: video.Id})
				if err != nil {
					logx.WithContext(l.ctx).Errorf("get video like count by like rpc failed, err: %v", err)
				}
				resp.Videos[index].CommentCount = countResp.Count
				resp.Videos[index].LikeCount = likeCountResp.LikeCount
			}(index, currentVideo)
		}
		wg.Wait()
	}
	return resp, nil
}
