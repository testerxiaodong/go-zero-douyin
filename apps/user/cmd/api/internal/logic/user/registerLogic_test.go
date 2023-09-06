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

func TestRegisterLogic_Register(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	// user rpc client的mock
	mockUserRpcClient := mock.NewMockUser(ctl)

	// validator mock
	mockValidator := globalMock.NewMockValidator(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, UserRpc: mockUserRpcClient}

	registerLogic := user.NewRegisterLogic(context.Background(), serviceContext)

	// 期望的结果
	expectedRegisterResp := &pb.RegisterUserResp{
		AccessToken:  utils.NewRandomString(10),
		RefreshAfter: utils.NewRandomInt64(1, 100),
		ExpireTime:   utils.NewRandomInt64(1, 100),
	}

	// 参数校验失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return(utils.NewRandomString(10))

	// 参数校验成功，但userRpc.GetUserInfo失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpcClient.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("userRpc.GetUserInfo error"))

	// 参数校验成功，且userRpc.GetUserInfo成功
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpcClient.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).
		Return(expectedRegisterResp, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *types.RegisterReq
	}{
		{
			name: "register_user_with_validator_error",
			req:  &types.RegisterReq{Username: "test", Password: "test"},
		},
		{
			name: "register_user_with_rpc_error",
			req:  &types.RegisterReq{Username: "test", Password: "test"},
		},
		{
			name: "register_user_success",
			req:  &types.RegisterReq{Username: "test", Password: "test"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := registerLogic.Register(tc.req)
			if err == nil {
				assert.NotEmpty(t, resp)
				assert.Equal(t, resp.AccessToken, expectedRegisterResp.GetAccessToken())
				assert.Equal(t, resp.RefreshAfter, expectedRegisterResp.GetRefreshAfter())
				assert.Equal(t, resp.ExpireTime, expectedRegisterResp.GetExpireTime())
			} else {
				assert.Nil(t, resp)
			}
		})
	}

}
