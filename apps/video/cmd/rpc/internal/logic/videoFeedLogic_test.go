package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/video/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/mock"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"testing"
)

func TestVideoFeedLogic_VideoFeed(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockVideoDo := mock.NewMockVideoDo(ctl)

	serviceContext := &svc.ServiceContext{VideoDo: mockVideoDo}

	videoFeedLogic := logic.NewVideoFeedLogic(context.Background(), serviceContext)

	// 查询数据库失败mock
	mockVideoDo.EXPECT().GetVideoListByTimeStampAndSectionId(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("search database error"))

	// 数据库没有数据mock
	mockVideoDo.EXPECT().GetVideoListByTimeStampAndSectionId(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*model.Video{}, nil)

	// 查询数据库成功mock
	mockVideoDo.EXPECT().GetVideoListByTimeStampAndSectionId(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*model.Video{NewRandomVideo(), NewRandomVideo()}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.VideoFeedReq
	}{
		{
			name: "get_video_feed_with_empty_param",
			req:  nil,
		},
		{
			name: "get_video_feed_with_empty_timestamp",
			req:  &pb.VideoFeedReq{LastTimeStamp: nil},
		},
		{
			name: "get_video_feed_with_search_database_error",
			req:  &pb.VideoFeedReq{LastTimeStamp: utils.FromInt64TimeStampToProtobufTimeStamp(100)},
		},
		{
			name: "get_video_feed_with_no_database_record",
			req:  &pb.VideoFeedReq{LastTimeStamp: utils.FromInt64TimeStampToProtobufTimeStamp(100)},
		},
		{
			name: "get_video_feed_success",
			req:  &pb.VideoFeedReq{LastTimeStamp: utils.FromInt64TimeStampToProtobufTimeStamp(100)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := videoFeedLogic.VideoFeed(tc.req)
			if err == nil {
				assert.NotNil(t, resp)
				if len(resp.GetVideos()) > 0 {
					assert.Equal(t, len(resp.GetVideos()), 2)
				}
			} else {
				assert.Empty(t, resp)
			}
		})
	}
}
