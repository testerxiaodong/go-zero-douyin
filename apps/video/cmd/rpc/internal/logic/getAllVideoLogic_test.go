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
	pbVideo "go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"
	"testing"
)

func TestGetAllVideoLogic_GetAllVideo(t *testing.T) {
	ctl := gomock.NewController(t)

	mockVideoDo := mock.NewMockVideoDo(ctl)

	serviceContext := &svc.ServiceContext{VideoDo: mockVideoDo}

	getAllVideoLogic := logic.NewGetAllVideoLogic(context.Background(), serviceContext)

	// VideoDo.GetAllVideo失败mock
	dbError := errors.New("VideoDo.GetAllVideo error")
	mockVideoDo.EXPECT().GetAllVideo(gomock.Any()).Return(nil, dbError)

	// 查询数据库成功，但没有数据的mock
	mockVideoDo.EXPECT().GetAllVideo(gomock.Any()).Return([]*model.Video{}, nil)

	// 查询数据库成功，有数据的mock
	expectedValue := NewRandomVideo()
	mockVideoDo.EXPECT().GetAllVideo(gomock.Any()).Return([]*model.Video{expectedValue, expectedValue}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "get_all_video_with_database_error",
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "查询数据库所有视频数据失败, err: %v", dbError),
		},
		{
			name: "get_all_video_with_no_record",
			err:  nil,
		},
		{
			name: "get_all_video_success",
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getAllVideoLogic.GetAllVideo(&pbVideo.GetAllVideoReq{})
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
