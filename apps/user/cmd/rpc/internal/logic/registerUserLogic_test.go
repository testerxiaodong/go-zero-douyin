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
	"gorm.io/gorm"
	"testing"
)

func TestRegisterUserLogic_RegisterUser(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockUserDo := mock.NewMockUserDo(ctl)

	serviceContext := &svc.ServiceContext{UserDo: mockUserDo, Config: config.Config{JwtAuth: struct {
		AccessSecret string
		AccessExpire int64
	}{AccessSecret: "test", AccessExpire: 600}}}

	registerUserLogic := logic.NewRegisterUserLogic(context.Background(), serviceContext)

	// 查询用户是否存在时失败mock
	mockUserDo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Return(nil, errors.New("search database error"))

	// 用户已存在mock
	mockUserDo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Return(&model.User{Username: "test", Password: "test", ID: 1}, nil)

	// 插入用户失败mock
	mockUserDo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	mockUserDo.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(errors.New("insert database error"))

	// 插入用户成功mock

	mockUserDo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	mockUserDo.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.RegisterUserReq
	}{
		{
			name: "login_with_empty_param",
			req:  nil,
		},
		{
			name: "login_with_empty_username",
			req:  &pb.RegisterUserReq{Password: "test"},
		},
		{
			name: "login_with_empty_password",
			req:  &pb.RegisterUserReq{Username: "test"},
		},
		{
			name: "login_with_database_search_error",
			req:  &pb.RegisterUserReq{Username: "test", Password: "test"},
		},
		{
			name: "login_with_had_database_record",
			req:  &pb.RegisterUserReq{Username: "test", Password: "test"},
		},
		{
			name: "login_with_insert_database_error",
			req:  &pb.RegisterUserReq{Username: "test", Password: "test"},
		},
		{
			name: "login_success",
			req:  &pb.RegisterUserReq{Username: "test", Password: "test"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			infoResp, err := registerUserLogic.RegisterUser(tc.req)
			if err == nil {
				assert.NotEmpty(t, infoResp)
			} else {
				assert.Empty(t, infoResp)
			}
		})
	}
}
