package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/search/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/search/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/search/cmd/rpc/pb"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"
	"go-zero-douyin/mock"
	"testing"
)

func TestCompleteVideoLogic_CompleteVideo(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockElasticsearch := mock.NewMockElasticService(ctl)
	serviceContext := &svc.ServiceContext{ElasticSearch: mockElasticsearch}
	completeVideoLogic := logic.NewCompleteVideoLogic(context.Background(), serviceContext)

	// 调用es失败的mock
	esError := errors.New("ElasticSearch.Suggestion error")
	mockElasticsearch.EXPECT().Suggestion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, esError)

	// 调用es成功的mock
	resp := &elastic.SearchResult{Suggest: map[string][]elastic.SearchSuggestion{
		xconst.ElasticSearchVideoSuggestionName: {{Text: "test", Options: []elastic.SearchSuggestionOption{{Text: "test"}}}},
	}}
	mockElasticsearch.EXPECT().Suggestion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(resp, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.CompleteVideoReq
		err  error
	}{
		{
			name: "complete_video_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "视频搜索自动补全的参数不能为空"),
		},
		{
			name: "complete_video_with_empty_input",
			req:  &pb.CompleteVideoReq{Input: ""},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "视频搜索自动补全的参数不能为空"),
		},
		{
			name: "complete_video_with_es_error",
			req:  &pb.CompleteVideoReq{Input: "test"},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "调用es自动补全失败, err: %v", esError),
		},
		{
			name: "complete_video_success",
			req:  &pb.CompleteVideoReq{Input: "test"},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := completeVideoLogic.CompleteVideo(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
