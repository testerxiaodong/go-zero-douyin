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
	"gorm.io/gen"
	"gorm.io/gorm"
	"testing"
)

func TestDeleteVideoLogic_DeleteVideo(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockVideoDo := mock.NewMockVideoDo(ctl)

	serviceContext := &svc.ServiceContext{VideoDo: mockVideoDo}

	deleteVideoLogic := logic.NewDeleteVideoLogic(context.Background(), serviceContext)

	// 通过id查询视频失败mock
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(nil, errors.New("search database error"))

	// 视频不存在mock
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	// 视频非该用户发布mock
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(&model.Video{ID: 1, OwnerID: 1}, nil)

	// 视频存在，且是该用户发布，但删除视频失败mock
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(&model.Video{ID: 1, OwnerID: 2}, nil)
	mockVideoDo.EXPECT().DeleteVideo(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, errors.New("delete database error"))

	// 视频存在，且是该用户发布，且视频删除成功mock
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(&model.Video{ID: 1, OwnerID: 2}, nil)
	mockVideoDo.EXPECT().DeleteVideo(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.DeleteVideoReq
	}{
		{
			name: "delete_video_with_empty_param",
			req:  nil,
		},
		{
			name: "delete_video_with_empty_user_id",
			req:  &pb.DeleteVideoReq{VideoId: 2},
		},
		{
			name: "delete_video_with_empty_video_id",
			req:  &pb.DeleteVideoReq{UserId: 2},
		},
		{
			name: "delete_video_with_database_search_error",
			req:  &pb.DeleteVideoReq{VideoId: 2, UserId: 2},
		},
		{
			name: "delete_video_with_no_database_record",
			req:  &pb.DeleteVideoReq{VideoId: 2, UserId: 2},
		},
		{
			name: "delete_video_with_is_not_video_owner",
			req:  &pb.DeleteVideoReq{VideoId: 2, UserId: 2},
		},
		{
			name: "delete_video_with_delete_video_error",
			req:  &pb.DeleteVideoReq{VideoId: 2, UserId: 2},
		},
		{
			name: "delete_video_with_delete_video_success",
			req:  &pb.DeleteVideoReq{VideoId: 2, UserId: 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			infoResp, err := deleteVideoLogic.DeleteVideo(tc.req)
			if err == nil {
				assert.NotNil(t, infoResp)
			} else {
				assert.Empty(t, infoResp)
			}
		})
	}
}
