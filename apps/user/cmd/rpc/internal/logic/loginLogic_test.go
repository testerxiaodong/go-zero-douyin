package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/user/cmd/rpc/internal/config"
	"go-zero-douyin/apps/user/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/user/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/user/cmd/rpc/mock"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	"gorm.io/gorm"
	"testing"
)

func TestLoginLogic_Login(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockUserDo := mock.NewMockUserDo(ctl)

	serviceContext := &svc.ServiceContext{UserDo: mockUserDo, Config: config.Config{JwtAuth: struct {
		AccessSecret string
		AccessExpire int64
	}{AccessSecret: "test", AccessExpire: 600}}}

	loginLogic := logic.NewLoginLogic(context.Background(), serviceContext)

	mockUserDo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Return(nil, errors.New("get user by username failed"))
	mockUserDo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockUserDo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Return(nil, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.LoginReq
	}{
		{
			name: "login_with_empty_param",
			req:  nil,
		},
		{
			name: "login_with_database_search_error",
			req:  &pb.LoginReq{Username: "test", Password: "test"},
		},
		{
			name: "login_with_no_database_record",
			req:  &pb.LoginReq{Username: "test", Password: "test"},
		},
		{
			name: "login_success",
			req:  &pb.LoginReq{Username: "test", Password: "test"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			infoResp, err := loginLogic.Login(tc.req)
			if err == nil {
				assert.NotEmpty(t, infoResp)
			} else {
				assert.Empty(t, infoResp)
			}
		})
	}

}
