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

func TestUserVideoListLogic_UserVideoList(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockVideoDo := mock.NewMockVideoDo(ctl)

	serviceContext := &svc.ServiceContext{VideoDo: mockVideoDo}

	userVideoListLogic := logic.NewUserVideoListLogic(context.Background(), serviceContext)

	// 查询数据库失败mock
	mockVideoDo.EXPECT().GetVideoListByUserId(gomock.Any(), gomock.Any()).Return(nil, errors.New("search database error"))

	// 数据库没有数据mock
	mockVideoDo.EXPECT().GetVideoListByUserId(gomock.Any(), gomock.Any()).Return([]*model.Video{}, nil)

	// 查询数据库成功
	mockVideoDo.EXPECT().GetVideoListByUserId(gomock.Any(), gomock.Any()).Return([]*model.Video{NewRandomVideo(), NewRandomVideo()}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.UserVideoListReq
	}{
		{
			name: "publish_video_with_empty_param",
			req:  nil,
		},
		{
			name: "publish_video_with_search_database_error",
			req:  &pb.UserVideoListReq{UserId: utils.NewRandomInt64(1, 10)},
		},
		{
			name: "publish_video_with_no_database_record",
			req:  &pb.UserVideoListReq{UserId: utils.NewRandomInt64(1, 10)},
		},
		{
			name: "publish_video_success",
			req:  &pb.UserVideoListReq{UserId: utils.NewRandomInt64(1, 10)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := userVideoListLogic.UserVideoList(tc.req)
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
