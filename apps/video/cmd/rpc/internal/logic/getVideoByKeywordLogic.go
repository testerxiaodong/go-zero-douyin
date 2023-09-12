package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"
	"reflect"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoByKeywordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoByKeywordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoByKeywordLogic {
	return &GetVideoByKeywordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoByKeywordLogic) GetVideoByKeyword(in *pb.GetVideoByKeywordReq) (*pb.GetVideoByKeywordResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if len(in.GetKeyword()) == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "查询信息不能为空")
	}

	// 调用es
	elasticsearchResp, err := l.svcCtx.Elasticsearch.SearchByKeyword(l.ctx, "video", in.GetKeyword())
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "通过关键字查询es数据失败, err: %v, req: %v", err, in)
	}
	fmt.Printf("查询消耗时间 %d ms, 结果总数: %d\n", elasticsearchResp.TookInMillis, elasticsearchResp.TotalHits())
	// 拼接响应
	if elasticsearchResp.TotalHits() > 0 {
		resp := &pb.GetVideoByKeywordResp{Videos: make([]*pb.VideoDetailInfo, 0, elasticsearchResp.TotalHits())}
		// 查询结果不为空，则遍历结果
		var video *pb.VideoDetailInfo
		// 通过Each方法，将es结果的json结构转换成struct对象
		for _, item := range elasticsearchResp.Each(reflect.TypeOf(video)) {
			// 转换成Article对象
			if v, ok := item.(*pb.VideoDetailInfo); ok {
				videoInfo := &pb.VideoDetailInfo{}
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
	return &pb.GetVideoByKeywordResp{}, nil
}
