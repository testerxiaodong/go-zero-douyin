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
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"gorm.io/gen"
	"gorm.io/gorm"
	"testing"
)

func TestDeleteVideoLogic_DeleteVideo(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockVideoDo := mock.NewMockVideoDo(ctl)
	mockSender := globalMock.NewMockSender(ctl)

	serviceContext := &svc.ServiceContext{VideoDo: mockVideoDo, Rabbit: mockSender}

	deleteVideoLogic := logic.NewDeleteVideoLogic(context.Background(), serviceContext)

	// 通过id查询视频失败mock
	dbSearchError := errors.New("search database error")
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 视频不存在mock
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	// 视频非该用户发布mock
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(&model.Video{ID: 1, OwnerID: 1}, nil)

	// 视频存在，且是该用户发布，但删除视频失败mock
	dbDeleteError := errors.New("delete database error")
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(&model.Video{ID: 1, OwnerID: 2}, nil)
	mockVideoDo.EXPECT().DeleteVideo(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, dbDeleteError)

	// 视频存在，且是该用户发布，且视频删除成功mock，但发布消息失败的mock
	senderError := errors.New("sender error")
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(&model.Video{ID: 1, OwnerID: 2}, nil)
	mockVideoDo.EXPECT().DeleteVideo(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 成功的mock
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(&model.Video{ID: 1, OwnerID: 2}, nil)
	mockVideoDo.EXPECT().DeleteVideo(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.DeleteVideoReq
		err  error
	}{
		{
			name: "delete_video_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Delete video with empty param"),
		},
		{
			name: "delete_video_with_empty_user_id",
			req:  &pb.DeleteVideoReq{VideoId: 2},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Delete video with empty user_id or video_id"),
		},
		{
			name: "delete_video_with_empty_video_id",
			req:  &pb.DeleteVideoReq{UserId: 2},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Delete video with empty user_id or video_id"),
		},
		{
			name: "delete_video_with_database_search_error",
			req:  &pb.DeleteVideoReq{VideoId: 2, UserId: 2},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "Find video by id failed: %v", dbSearchError),
		},
		{
			name: "delete_video_with_no_database_record",
			req:  &pb.DeleteVideoReq{VideoId: 2, UserId: 2},
			err:  errors.Wrapf(logic.ErrVideoNotFound, "video_id: %d", 2),
		},
		{
			name: "delete_video_with_is_not_video_owner",
			req:  &pb.DeleteVideoReq{VideoId: 2, UserId: 2},
			err:  errors.Wrapf(xerr.NewErrMsg("视频非该用户发布，用户无权操作"), "video_id: %d", 2),
		},
		{
			name: "delete_video_with_delete_video_error",
			req:  &pb.DeleteVideoReq{VideoId: 2, UserId: 2},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "delete video by id failed: %v", dbDeleteError),
		},
		{
			name: "delete_video_with_sender_error",
			req:  &pb.DeleteVideoReq{VideoId: 2, UserId: 2},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.RPC_DELETE_ERR), "req: %v, err: %v", &pb.DeleteVideoReq{VideoId: 2, UserId: 2}, senderError),
		},
		{
			name: "delete_video_with_delete_video_success",
			req:  &pb.DeleteVideoReq{VideoId: 2, UserId: 2},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := deleteVideoLogic.DeleteVideo(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
