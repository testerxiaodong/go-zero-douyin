package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/search/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/search/cmd/rpc/pb"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserLogic {
	return &SearchUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUserLogic) SearchUser(in *pb.SearchUserReq) (*pb.SearchUserResp, error) {
	// todo: add your logic here and delete this line
	// 参数业务逻辑处理
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "参数不能为nil")
	}
	if len(in.GetKeyword()) == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "用户搜索关键字不能为空")
	}
	// 搜索业务逻辑处理
	var sort string
	if in.GetSort() == 1 {
		sort = "follower_count"
	} else {
		sort = "_score"
	}
	searchResult, err := l.svcCtx.ElasticSearch.SearchByKeyword(l.ctx, xconst.ElasticSearchUserIndexName, "username", in.GetKeyword(), in.GetPage(), in.GetPageSize(), sort, in.GetHighlight())
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "根据关键字搜索用户信息失败, err: %v, req: %v", err, in)
	}
	// 没有数据，直接返回
	if searchResult.TotalHits() == 0 {
		return &pb.SearchUserResp{}, nil
	}
	// 有数据，拼接响应
	resp := &pb.SearchUserResp{Users: make([]*pb.User, 0, searchResult.TotalHits())}
	resp.Total = searchResult.TotalHits()
	// 查询结果不为空，则遍历结果
	for _, item := range searchResult.Hits.Hits {
		var user pb.User
		// 反序列化source字段
		if err := json.Unmarshal(item.Source, &user); err != nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("反序列化es搜索结果失败"), "err: %v", err)
		}
		// 如果有高亮需求，替换高亮字段文本
		if item.Highlight != nil {
			user.Username = item.Highlight["username"][0]
		}
		resp.Users = append(resp.Users, &user)
	}
	return resp, nil
}
