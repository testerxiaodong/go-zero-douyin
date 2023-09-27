package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/recommend/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/recommend/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoRecommendSectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoRecommendSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoRecommendSectionLogic {
	return &VideoRecommendSectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VideoRecommendSectionLogic) VideoRecommendSection(in *pb.VideoRecommendSectionReq) (*pb.VideoRecommendSectionResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "获取分区视频推荐时参数为nil")
	}
	if in.GetSectionId() == 0 || in.GetUserId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "获取分区视频推荐时分区id或者用户id不能为空")
	}
	// 调用GorseApi
	var videoIds []string
	queryParams := map[string]string{"n": cast.ToString(in.GetCount()), "write-back-type": "read", "write-back-delay": "10m"}
	url := fmt.Sprintf("%s/api/recommend/%d/%d", l.svcCtx.GorseConf, in.GetUserId(), in.GetSectionId())
	err := l.svcCtx.RestClient.Get(queryParams, url, &videoIds)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "GorseApi获取推荐失败, err: %v, req: %v", err, in)
	}
	// 拼接响应
	fmt.Println(videoIds)
	if len(videoIds) > 0 {
		resp := &pb.VideoRecommendSectionResp{VideoIds: make([]int64, 0)}
		for _, idStr := range videoIds {
			resp.VideoIds = append(resp.VideoIds, cast.ToInt64(idStr))
		}
		return resp, nil
	}
	return &pb.VideoRecommendSectionResp{}, nil
}
