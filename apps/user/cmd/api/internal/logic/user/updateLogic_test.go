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
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestUpdateLogic_Update(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	// user rpc client的mock
	mockUserRpcClient := mock.NewMockUser(ctl)

	// validator mock
	mockValidator := globalMock.NewMockValidator(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, UserRpc: mockUserRpcClient}

	updateLogic := user.NewUpdateLogic(context.Background(), serviceContext)

	validateResult := utils.NewRandomString(10)

	// 参数校验失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// 参数校验成功，但userRpc.UpdateUser失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")

	rpcError := errors.New("userRpc.UpdateUser error")

	mockUserRpcClient.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).
		Return(nil, rpcError)

	// 参数校验成功，且userRpc.GetUserInfo成功
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpcClient.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).
		Return(&pb.UpdateUserResp{}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *types.UpdateUserReq
		err  error
	}{
		{
			name: "update_user_with_validator_error",
			req:  &types.UpdateUserReq{Username: "test", Password: "test"},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "update_user_with_rpc_error",
			req:  &types.UpdateUserReq{Username: "test", Password: "test"},
			err:  errors.Wrapf(rpcError, "req: %v", &types.UpdateUserReq{Username: "test", Password: "test"}),
		},
		{
			name: "update_user_success",
			req:  &types.UpdateUserReq{Username: "test", Password: "test"},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := updateLogic.Update(tc.req)
			if err == nil {
				assert.Equal(t, err, tc.err)
			} else {
				assert.Equal(t, err.Error(), tc.err.Error())
			}
		})
	}
}
