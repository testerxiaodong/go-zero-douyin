package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"
	"reflect"

	"go-zero-douyin/apps/search/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/search/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchVideoLogic {
	return &SearchVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchVideoLogic) SearchVideo(in *pb.SearchVideoReq) (*pb.SearchVideoResp, error) {
	// todo: add your logic here and delete this line
	// 参数业务逻辑处理
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "参数不能为nil")
	}
	if len(in.GetKeyword()) == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "视频搜索关键字不能为空")
	}

	var sort string
	if in.GetSort() == 1 {
		sort = "like_count"
	} else if in.GetSort() == 2 {
		sort = "comment_count"
	} else {
		sort = "create_time"
	}
	searchResult, err := l.svcCtx.ElasticSearch.SearchByKeyword(l.ctx, xconst.ElasticSearchVideoIndexName, "title", in.GetKeyword(), in.GetPage(), in.GetPageSize(), sort)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "根据关键字搜索视频信息失败, err: %v, req: %v", err, in)
	}

	if searchResult.TotalHits() > 0 {
		resp := &pb.SearchVideoResp{Videos: make([]*pb.Video, 0, searchResult.TotalHits())}
		resp.Total = searchResult.TotalHits()
		// 查询结果不为空，则遍历结果
		var video *pb.Video
		// 通过Each方法，将es结果的json结构转换成struct对象
		for _, item := range searchResult.Each(reflect.TypeOf(video)) {
			// 转换成Article对象
			if v, ok := item.(*pb.Video); ok {
				videoInfo := &pb.Video{}
				videoInfo.Id = v.Id
				videoInfo.Title = v.Title
				videoInfo.SectionId = v.SectionId
				videoInfo.Tags = v.Tags
				videoInfo.OwnerId = v.OwnerId
				videoInfo.PlayUrl = v.PlayUrl
				videoInfo.CoverUrl = v.CoverUrl
				videoInfo.CommentCount = v.CommentCount
				videoInfo.LikeCount = v.LikeCount
				videoInfo.CreateTime = v.CreateTime
				videoInfo.UpdateTime = v.UpdateTime
				resp.Videos = append(resp.Videos, videoInfo)
			}
		}
		return resp, nil
	}
	return &pb.SearchVideoResp{}, nil
}
