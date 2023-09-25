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

func TestUpdateUserLogic_UpdateUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserModel := mock.NewMockuserModel(ctl)
	serviceContext := &svc.ServiceContext{UserModel: mockUserModel}
	updateUserLogic := logic.NewUpdateUserLogic(context.Background(), serviceContext)

	// 根据id查询用户信息失败的mock
	dbError := errors.New("search database error")
	mockUserModel.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, dbError)

	// 根据id查询用户不存在
	mockUserModel.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, model.ErrNotFound)

	// 查询用户信息失败mock
	mockUserModel.EXPECT().FindOne(gomock.Any(), gomock.Any()).
		Return(&model.User{Id: 1, Username: "test", Password: "test"}, nil)
	mockUserModel.EXPECT().FindOneByUsernameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, dbError)

	// 用户名已存在mock
	mockUserModel.EXPECT().FindOne(gomock.Any(), gomock.Any()).
		Return(&model.User{Id: 1, Username: "test", Password: "test"}, nil)
	mockUserModel.EXPECT().FindOneByUsernameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.User{Id: 2}, nil)

	// 更新用户失败mock
	dbUpdateError := errors.New("update database error")
	mockUserModel.EXPECT().FindOne(gomock.Any(), gomock.Any()).
		Return(&model.User{Id: 1, Username: "test", Password: "test"}, nil)
	mockUserModel.EXPECT().FindOneByUsernameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.User{Username: "test", Password: "test", Id: 1}, nil)
	mockUserModel.EXPECT().UpdateWithVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbUpdateError)

	// 更新用户成功的mock
	mockUserModel.EXPECT().FindOne(gomock.Any(), gomock.Any()).
		Return(&model.User{Id: 1, Username: "test", Password: "test"}, nil)
	mockUserModel.EXPECT().FindOneByUsernameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.User{Username: "test", Password: "test", Id: 1}, nil)
	mockUserModel.EXPECT().UpdateWithVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.UpdateUserReq
		err  error
	}{
		{
			name: "update_user_with_empty_param",
			req:  nil,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Update user empty param"),
		},
		{
			name: "update_user_with_empty_id",
			req:  &pb.UpdateUserReq{Password: "test", Username: "test"},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Updtae user empty user id"),
		},
		{
			name: "update_user_with_search_database_error_by_id",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"根据id查询用户失败, err: %v, user_id: %d", dbError, 1),
		},
		{
			name: "update_user_with_search_database_no_record_by_id",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
			err:  errors.Wrapf(logic.ErrUserNotFound, "user_id: %d", 1),
		},
		{
			name: "update_user_with_search_database_error",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"根据用户名查询用户失败, err: %v, username: %s", dbError, "test"),
		},
		{
			name: "update_user_with_username_exist",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
			err:  errors.Wrapf(logic.ErrUserAlreadyRegister, "username: %s", "test"),
		},
		{
			name: "update_user_with_update_database_error",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR), "update mysql user info failed, err: %v", dbUpdateError),
		},
		{
			name: "update_success",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := updateUserLogic.UpdateUser(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
