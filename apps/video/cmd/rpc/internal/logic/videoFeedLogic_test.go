package logic_test

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/video/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/mock"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"
	"testing"
)

func TestVideoFeedLogic_VideoFeed(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockVideoDo := mock.NewMockvideoModel(ctl)
	serviceContext := &svc.ServiceContext{VideoModel: mockVideoDo}
	videoFeedLogic := logic.NewVideoFeedLogic(context.Background(), serviceContext)

	// 查询数据库失败mock
	dbError := errors.New("search database error")
	mockVideoDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockVideoDo.EXPECT().FindPageListByPage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, dbError)

	// 数据库没有数据mock
	mockVideoDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockVideoDo.EXPECT().FindPageListByPage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*model.Video{}, nil)

	// 查询数据库成功mock
	mockVideoDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockVideoDo.EXPECT().FindPageListByPage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*model.Video{NewRandomVideo(), NewRandomVideo()}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.VideoFeedReq
		err  error
	}{
		{
			name: "get_video_feed_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Get video feed with empty param"),
		},
		{
			name: "get_video_feed_with_empty_timestamp",
			req:  &pb.VideoFeedReq{LastTimeStamp: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Get video feed with empty timestamp"),
		},
		{
			name: "get_video_feed_with_search_database_error",
			req:  &pb.VideoFeedReq{LastTimeStamp: 10},
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"Get video feed by last timestmap failed, err: %v", dbError),
		},
		{
			name: "get_video_feed_with_no_database_record",
			req:  &pb.VideoFeedReq{LastTimeStamp: 10},
			err:  nil,
		},
		{
			name: "get_video_feed_success",
			req:  &pb.VideoFeedReq{LastTimeStamp: 10},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := videoFeedLogic.VideoFeed(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
