package video

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

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
	feedResp, err := l.svcCtx.VideoRpc.VideoFeed(l.ctx, &pb.VideoFeedReq{LastTimeStamp: utils.FromInt64TimeStampToProtobufTimeStamp(req.LastTimeStamp)})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}
	resp = &types.VideoFeedResp{Videos: make([]*types.Video, 0)}
	err = copier.Copy(resp, feedResp)
	if err != nil {
		return nil, errors.Wrapf(err, "copier feed resp failed: %v", feedResp)
	}
	return resp, nil
}
