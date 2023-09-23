package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/user/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/user/cmd/rpc/internal/model"
	"go-zero-douyin/apps/user/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/user/cmd/rpc/mock"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"
	"testing"
)

func TestGetUserInfoLogic_GetUserInfo(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserModel := mock.NewMockuserModel(ctl)
	serviceContext := &svc.ServiceContext{UserModel: mockUserModel}
	getUserInfoLogic := logic.NewGetUserInfoLogic(context.Background(), serviceContext)
	dbError := errors.New("database search error")
	mockUserModel.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, dbError)
	mockUserModel.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, model.ErrNotFound)
	mockUserModel.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(&model.User{Id: 2, Username: "test"}, nil)
	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.GetUserInfoReq
		err  error
	}{
		{
			name: "get_user_info_with_empty_param",
			req:  nil,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "获取用户信息时参数为nil"),
		},
		{
			name: "get_user_info_with_empty_id",
			req:  &pb.GetUserInfoReq{Id: 0},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "获取用户信息时id为空"),
		},
		{
			name: "get_user_info_with_database_search_error",
			req:  &pb.GetUserInfoReq{Id: 2},
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"从数据库查询用户信息失败 user_id:%d , err:%v", 2, dbError),
		},
		{
			name: "get_user_info_with_no_database_record",
			req:  &pb.GetUserInfoReq{Id: 2},
			err:  errors.Wrapf(logic.ErrUserNotFound, "id:%d", 2),
		},
		{
			name: "get_user_info_success",
			req:  &pb.GetUserInfoReq{Id: 2},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getUserInfoLogic.GetUserInfo(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
