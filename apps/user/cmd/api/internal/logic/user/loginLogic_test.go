package user_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/user/cmd/api/internal/logic/user"
	"go-zero-douyin/apps/user/cmd/api/internal/svc"
	"go-zero-douyin/apps/user/cmd/api/internal/types"
	"go-zero-douyin/apps/user/cmd/rpc/mock"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestLoginLogic_Login(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	// user rpc client的mock
	mockUserRpcClient := mock.NewMockUser(ctl)

	// validator mock
	mockValidator := globalMock.NewMockValidator(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, UserRpc: mockUserRpcClient}

	loginLogic := user.NewLoginLogic(context.Background(), serviceContext)

	// 期望的结果
	expectedLoginResp := &pb.LoginResp{
		AccessToken:  utils.NewRandomString(10),
		RefreshAfter: utils.NewRandomInt64(1, 10),
		ExpireTime:   utils.NewRandomInt64(1, 100),
	}

	// 参数校验失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return(utils.NewRandomString(10))

	// 参数校验成功，但userRpc.Login失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpcClient.EXPECT().Login(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("userRpc.Login error"))

	// 参数校验成功，且userRpc.Login成功的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpcClient.EXPECT().Login(gomock.Any(), gomock.Any()).
		Return(expectedLoginResp, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *types.LoginReq
	}{
		{
			name: "login_with_validator_error",
			req:  &types.LoginReq{Username: "test", Password: "test"},
		},
		{
			name: "login_with_rpc_error",
			req:  &types.LoginReq{Username: "test", Password: "test"},
		},
		{
			name: "login_success",
			req:  &types.LoginReq{Username: "test", Password: "test"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := loginLogic.Login(tc.req)
			if err == nil {
				assert.NotEmpty(t, resp)
				assert.Equal(t, resp.AccessToken, expectedLoginResp.GetAccessToken())
				assert.Equal(t, resp.RefreshAfter, expectedLoginResp.GetRefreshAfter())
				assert.Equal(t, resp.ExpireTime, expectedLoginResp.GetExpireTime())
			} else {
				assert.Nil(t, resp)
			}
		})
	}
}
