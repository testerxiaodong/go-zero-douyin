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
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"testing"
)

func TestUserVideoListLogic_UserVideoList(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockVideoDo := mock.NewMockvideoModel(ctl)

	serviceContext := &svc.ServiceContext{VideoModel: mockVideoDo}

	userVideoListLogic := logic.NewUserVideoListLogic(context.Background(), serviceContext)

	// 查询数据库失败mock
	dbError := errors.New("search database error")
	mockVideoDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockVideoDo.EXPECT().FindPageListByPageWithTotal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, int64(0), dbError)

	// 数据库没有数据mock
	mockVideoDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockVideoDo.EXPECT().FindPageListByPageWithTotal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*model.Video{}, int64(0), nil)

	// 查询数据库成功
	mockVideoDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockVideoDo.EXPECT().FindPageListByPageWithTotal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*model.Video{NewRandomVideo(), NewRandomVideo()}, int64(2), nil)

	// 表格驱动测试
	req := &pb.UserVideoListReq{
		UserId:   utils.NewRandomInt64(1, 10),
		Page:     utils.NewRandomInt64(1, 10),
		PageSize: utils.NewRandomInt64(1, 10),
	}
	testCases := []struct {
		name string
		req  *pb.UserVideoListReq
		err  error
	}{
		{
			name: "get_user_video_list_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Get User Video List with empty param"),
		},
		{
			name: "get_user_video_list_with_empty_user_id",
			req:  &pb.UserVideoListReq{UserId: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Get user video list with empty user_id"),
		},
		{
			name: "get_user_video_list_with_search_database_error",
			req:  req,
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"Get user video list by user_id failed, err: %v", dbError),
		},
		{
			name: "get_user_video_list_with_no_database_record",
			req:  req,
			err:  nil,
		},
		{
			name: "get_user_video_list_success",
			req:  req,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := userVideoListLogic.UserVideoList(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
