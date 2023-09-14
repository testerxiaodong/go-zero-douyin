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

	var sort string
	if in.GetSort() == 1 {
		sort = "follower_count"
	} else {
		sort = "_score"
	}
	searchResult, err := l.svcCtx.ElasticSearch.SearchByKeyword(l.ctx, xconst.ElasticSearchUserIndexName, "username", in.GetKeyword(), in.GetPage(), in.GetPageSize(), sort)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "根据关键字搜索用户信息失败, err: %v, req: %v", err, in)
	}

	if searchResult.TotalHits() > 0 {
		resp := &pb.SearchUserResp{Users: make([]*pb.User, 0, searchResult.TotalHits())}
		resp.Total = searchResult.TotalHits()
		// 查询结果不为空，则遍历结果
		var video *pb.User
		// 通过Each方法，将es结果的json结构转换成struct对象
		for _, item := range searchResult.Each(reflect.TypeOf(video)) {
			// 转换成Article对象
			if v, ok := item.(*pb.User); ok {
				userInfo := &pb.User{}
				userInfo.Id = v.Id
				userInfo.Username = v.Username
				userInfo.FollowerCount = v.FollowerCount
				userInfo.FollowCount = v.FollowCount
				resp.Users = append(resp.Users, userInfo)
			}
		}
		return resp, nil
	}
	return &pb.SearchUserResp{}, nil
}
