package logic_test

import (
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-zero-douyin/apps/video/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestGetVideoByKeywordLogic_GetVideoByKeyword(t *testing.T) {
	ctl := gomock.NewController(t)

	mockElasticsearch := globalMock.NewMockElasticService(ctl)

	serviceContext := &svc.ServiceContext{Elasticsearch: mockElasticsearch}

	getVideoByKeywordLogic := logic.NewGetVideoByKeywordLogic(context.Background(), serviceContext)

	// es失败的mock
	esError := errors.New("Elasticsearch.SearchByKeyword error")
	mockElasticsearch.EXPECT().SearchByKeyword(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, esError)

	// es成功，但没有数据的mock
	mockElasticsearch.EXPECT().SearchByKeyword(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&elastic.SearchResult{}, nil)

	// es成功，有数据
	video := NewRandomVideo()
	data, err := json.Marshal(video)
	require.NoError(t, err)

	expectedValue := &elastic.SearchResult{Hits: &elastic.SearchHits{
		TotalHits: &elastic.TotalHits{Value: 1, Relation: "eq"},
		Hits:      []*elastic.SearchHit{&elastic.SearchHit{Source: data}},
	}}

	mockElasticsearch.EXPECT().SearchByKeyword(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(expectedValue, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.GetVideoByKeywordReq
		err  error
	}{
		{
			name: "get_video_by_keyword_with_empty_keyword",
			req:  &pb.GetVideoByKeywordReq{Keyword: ""},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "查询信息不能为空"),
		},
		{
			name: "get_video_by_keyword_with_es_error",
			req:  &pb.GetVideoByKeywordReq{Keyword: "test"},
			err: errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR),
				"通过关键字查询es数据失败, err: %v, req: %v", esError, &pb.GetVideoByKeywordReq{Keyword: "test"}),
		},
		{
			name: "get_video_by_keyword_with_no_record",
			req:  &pb.GetVideoByKeywordReq{Keyword: "test"},
			err:  nil,
		},
		{
			name: "get_video_by_keyword_success",
			req:  &pb.GetVideoByKeywordReq{Keyword: "test"},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getVideoByKeywordLogic.GetVideoByKeyword(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
