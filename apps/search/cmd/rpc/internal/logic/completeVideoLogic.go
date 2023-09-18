package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/search/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/search/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompleteVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompleteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteVideoLogic {
	return &CompleteVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CompleteVideoLogic) CompleteVideo(in *pb.CompleteVideoReq) (*pb.CompleteVideoResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil || len(in.GetInput()) == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "视频搜索自动补全的参数不能为空")
	}
	// 调用es
	result, err := l.svcCtx.ElasticSearch.Suggestion(l.ctx, xconst.ElasticSearchVideoIndexName, xconst.ElasticSearchVideoSuggestionName, in.GetInput())
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "调用es自动补全失败, err: %v", err)
	}
	// 拼接响应
	if result.Suggest != nil {
		suggestions := make([]string, 0)
		for _, suggest := range result.Suggest[xconst.ElasticSearchVideoSuggestionName] {
			for _, option := range suggest.Options {
				suggestions = append(suggestions, option.Text)
			}
		}
		return &pb.CompleteVideoResp{Suggestions: suggestions}, nil
	}

	return &pb.CompleteVideoResp{}, nil
}
