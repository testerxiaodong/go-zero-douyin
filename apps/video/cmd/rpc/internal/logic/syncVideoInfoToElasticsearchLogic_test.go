package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/video/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestSyncVideoInfoToElasticsearchLogic_SyncVideoInfoToElasticsearch(t *testing.T) {
	ctl := gomock.NewController(t)

	mockElasticsearch := globalMock.NewMockElasticService(ctl)

	serviceContext := &svc.ServiceContext{Elasticsearch: mockElasticsearch}

	syncVideoInfoToElasticsearchLogic := logic.NewSyncVideoInfoToElasticsearchLogic(context.Background(), serviceContext)

	// Elasticsearch.CreateDocument失败的mock
	esError := errors.New("Elasticsearch.CreateDocument error")
	mockElasticsearch.EXPECT().CreateDocument(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, esError)

	// Elasticsearch.CreateDocument成功的mock
	mockElasticsearch.EXPECT().CreateDocument(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&elastic.IndexResponse{}, nil)

	req := &pb.SyncVideoInfoToElasticsearchReq{VideoInfo: &pb.VideoDetailInfo{Id: 1}}

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.SyncVideoInfoToElasticsearchReq
		err  error
	}{
		{
			name: "sync_to_es_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "同步视频信息到es的参数为空"),
		},
		{
			name: "sync_to_es_with_es_error",
			req:  req,
			err: errors.Wrapf(xerr.NewErrCode(xerr.RPC_INSERT_ERR), "创建es文档失败, err: %v, video: %v",
				esError, req.GetVideoInfo()),
		},
		{
			name: "sync_to_es_success",
			req:  req,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := syncVideoInfoToElasticsearchLogic.SyncVideoInfoToElasticsearch(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
