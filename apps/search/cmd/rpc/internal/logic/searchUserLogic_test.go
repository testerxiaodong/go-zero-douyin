package logic_test

import (
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-zero-douyin/apps/search/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/search/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/search/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestSearchUserLogic_SearchUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockElastic := globalMock.NewMockElasticService(ctl)
	serviceContext := &svc.ServiceContext{ElasticSearch: mockElastic}
	searchUserLogic := logic.NewSearchUserLogic(context.Background(), serviceContext)

	// ElasticSearch.SearchByKeyword失败的mock
	esError := errors.New("es error")
	mockElastic.EXPECT().SearchByKeyword(gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, esError)

	// ElasticSearch.SearchByKeyword成功的mock
	user, err := json.Marshal(&pb.User{
		Id:            utils.NewRandomInt64(1, 10),
		Username:      utils.NewRandomString(10),
		FollowerCount: utils.NewRandomInt64(1, 10),
		FollowCount:   utils.NewRandomInt64(1, 10),
	})
	require.NoError(t, err)

	resp := &elastic.SearchResult{
		Hits: &elastic.SearchHits{
			TotalHits: &elastic.TotalHits{
				Value: 2,
			},
			Hits: []*elastic.SearchHit{
				{Source: user}, {Source: user},
			},
		},
	}
	mockElastic.EXPECT().SearchByKeyword(gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(resp, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.SearchUserReq
		err  error
	}{
		{
			name: "search_user_with_nil_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "参数不能为nil"),
		},
		{
			name: "search_user_with_empty_keyword",
			req:  &pb.SearchUserReq{Keyword: ""},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "用户搜索关键字不能为空"),
		},
		{
			name: "search_user_with_er_error",
			req:  &pb.SearchUserReq{Keyword: "test", Sort: 0},
			err: errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR),
				"根据关键字搜索用户信息失败, err: %v, req: %v", esError, &pb.SearchUserReq{Keyword: "test", Sort: 0}),
		},
		{
			name: "search_user_success",
			req:  &pb.SearchUserReq{Keyword: "test", Sort: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := searchUserLogic.SearchUser(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
