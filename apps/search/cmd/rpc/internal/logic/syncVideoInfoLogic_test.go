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
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestSyncVideoInfoLogic_SyncVideoInfo(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockElastic := globalMock.NewMockElasticService(ctl)
	serviceContext := &svc.ServiceContext{ElasticSearch: mockElastic}
	syncVideoInfoLogic := logic.NewSyncVideoInfoLogic(context.Background(), serviceContext)

	// ElasticSearch.DeleteDocument失败的mock
	esError := errors.New("es error")
	mockElastic.EXPECT().CreateDocument(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, esError)

	// ElasticSearch.DeleteDocument成功的mock
	mockElastic.EXPECT().CreateDocument(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&elastic.IndexResponse{}, nil)

	// 表格驱动测试
	req := &pb.SyncVideoInfoReq{Video: &pb.Video{
		Id:           utils.NewRandomInt64(1, 10),
		Title:        utils.NewRandomString(10),
		SectionId:    utils.NewRandomInt64(1, 10),
		Tags:         []string{"1", "2", "3"},
		OwnerId:      utils.NewRandomInt64(1, 10),
		PlayUrl:      utils.NewRandomString(10),
		CoverUrl:     utils.NewRandomString(10),
		CommentCount: utils.NewRandomInt64(1, 10),
		LikeCount:    utils.NewRandomInt64(1, 10),
		CreateTime:   utils.NewRandomInt64(1, 10),
		UpdateTime:   utils.NewRandomInt64(1, 10),
	}}
	testCases := []struct {
		name string
		req  *pb.SyncVideoInfoReq
		err  error
	}{
		{
			name: "sync_video_info_with_nil_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "参数不能为nil"),
		},
		{
			name: "sync_video_info_with_empty_id",
			req:  &pb.SyncVideoInfoReq{Video: &pb.Video{Id: 0}},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "视频id不允许为空"),
		},
		{
			name: "sync_video_info_with_es_error",
			req:  req,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.RPC_UPDATE_ERR), "更新es视频文档失败, err: %v", esError),
		},
		{
			name: "sync_video_info_success",
			req:  req,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := syncVideoInfoLogic.SyncVideoInfo(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
