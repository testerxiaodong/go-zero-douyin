package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/user/cmd/rpc/internal/config"
	"go-zero-douyin/apps/user/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/user/cmd/rpc/internal/model"
	"go-zero-douyin/apps/user/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/user/cmd/rpc/mock"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestRegisterUserLogic_RegisterUser(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockUserModel := mock.NewMockuserModel(ctl)
	serviceContext := &svc.ServiceContext{UserModel: mockUserModel, Config: config.Config{JwtAuth: struct {
		AccessSecret string
		AccessExpire int64
	}{AccessSecret: "test", AccessExpire: 600}}}

	registerUserLogic := logic.NewRegisterUserLogic(context.Background(), serviceContext)

	// 查询用户是否存在时失败mock
	dbSearchError := errors.New("search database error")
	mockUserModel.EXPECT().FindOneByUsernameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, dbSearchError)

	// 用户已存在mock
	mockUserModel.EXPECT().FindOneByUsernameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.User{Username: "test", Password: "test", Id: 1}, nil)

	// 插入用户失败mock
	dbInsertError := errors.New("insert database error")
	mockUserModel.EXPECT().FindOneByUsernameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, model.ErrNotFound)
	mockUserModel.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbInsertError)

	mockResult := globalMock.NewMockResult(ctl)
	// 成功mock
	mockUserModel.EXPECT().FindOneByUsernameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, model.ErrNotFound)
	mockUserModel.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil)
	mockResult.EXPECT().LastInsertId().Return(int64(10), nil)
	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.RegisterUserReq
		err  error
	}{
		{
			name: "register_with_empty_param",
			req:  nil,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Register user empty param"),
		},
		{
			name: "register_with_empty_username",
			req:  &pb.RegisterUserReq{Password: "test"},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Register user error param"),
		},
		{
			name: "register_with_empty_password",
			req:  &pb.RegisterUserReq{Username: "test"},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Register user error param"),
		},
		{
			name: "register_with_database_search_error",
			req:  &pb.RegisterUserReq{Username: "test", Password: "test"},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "find user by username failed, username: %s, err: %v", "test", dbSearchError),
		},
		{
			name: "register_with_had_database_record",
			req:  &pb.RegisterUserReq{Username: "test", Password: "test"},
			err:  errors.Wrapf(logic.ErrUserAlreadyRegister, "register user exists username: %s", "test"),
		},
		{
			name: "register_with_insert_database_error",
			req:  &pb.RegisterUserReq{Username: "test", Password: "test"},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "insert user failed, username: %s, password: %s", "test", "test"),
		},
		{
			name: "register_success",
			req:  &pb.RegisterUserReq{Username: "test", Password: "test"},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := registerUserLogic.RegisterUser(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
