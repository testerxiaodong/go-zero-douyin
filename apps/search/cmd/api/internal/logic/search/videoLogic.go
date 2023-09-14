package search

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/search/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/search/cmd/api/internal/svc"
	"go-zero-douyin/apps/search/cmd/api/internal/types"

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

func (l *VideoLogic) Video(req *types.SearchVideoReq) (resp *types.SearchVideoResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}

	var sort pb.SearchVideoReq_Sort
	if req.Sort == 1 {
		sort = pb.SearchVideoReq_LIKECOUNT
	} else if req.Sort == 2 {
		sort = pb.SearchVideoReq_COMMENTCOUNT
	} else {
		sort = pb.SearchVideoReq_DEFAULT
	}

	// 调用searchRpc
	searchVideoResp, err := l.svcCtx.SearchRpc.SearchVideo(l.ctx, &pb.SearchVideoReq{
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
		Sort:     sort,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}

	// 拼接响应
	if len(searchVideoResp.GetVideos()) > 0 {
		resp = &types.SearchVideoResp{Videos: make([]*types.Video, 0)}
		_ = copier.Copy(resp, searchVideoResp)
		return resp, nil
	}

	return &types.SearchVideoResp{}, nil
}
