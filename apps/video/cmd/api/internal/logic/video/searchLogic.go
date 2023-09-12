package video

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	pbVideo "go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchVideoReq) (resp *types.SearchVideoResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}

	// 调用videoRpc
	videos, err := l.svcCtx.VideoRpc.GetVideoByKeyword(l.ctx, &pbVideo.GetVideoByKeywordReq{Keyword: req.Keyword})
	if err != nil {
		return nil, errors.Wrapf(err, "VideoRpc.GetVideoByKeyword失败, req: %v", req)
	}

	// 拷贝响应
	if len(videos.GetVideos()) > 0 {
		resp = &types.SearchVideoResp{Videos: make([]*types.VideoInfo, 0)}
		err := copier.Copy(resp, videos)
		if err != nil {
			return nil, errors.Wrapf(err, "拷贝响应数据失败, data: %v", videos.GetVideos())
		}
		return resp, nil
	}

	return resp, nil
}
