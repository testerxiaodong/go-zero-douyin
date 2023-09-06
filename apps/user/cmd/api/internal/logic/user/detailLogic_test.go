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

func TestDetailLogic_Detail(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	// user rpc client的mock
	mockUserRpcClient := mock.NewMockUser(ctl)

	// validator mock
	mockValidator := globalMock.NewMockValidator(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, UserRpc: mockUserRpcClient}

	detailLogic := user.NewDetailLogic(context.Background(), serviceContext)

	// 期望的结果
	expectedUserInfo := &pb.UserInfo{
		Id:       utils.NewRandomInt64(1, 10),
		Username: utils.NewRandomString(10),
	}

	// 参数校验失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return(utils.NewRandomString(10))

	// 参数校验成功，但userRpc.GetUserInfo失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpcClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("userRpc.GetUserInfo error"))

	// 参数校验成功，且userRpc.GetUserInfo成功
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpcClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).
		Return(&pb.GetUserInfoResp{User: expectedUserInfo}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *types.UserInfoReq
	}{
		{
			name: "get_user_detail_with_validator_error",
			req:  &types.UserInfoReq{Id: 1},
		},
		{
			name: "get_user_detail_with_rpc_error",
			req:  &types.UserInfoReq{Id: 1},
		},
		{
			name: "get_user_detail_success",
			req:  &types.UserInfoReq{Id: 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := detailLogic.Detail(tc.req)
			if err == nil {
				assert.NotEmpty(t, resp)
				assert.Equal(t, resp.User.Id, expectedUserInfo.GetId())
				assert.Equal(t, resp.User.Username, expectedUserInfo.GetUsername())
			} else {
				assert.Nil(t, resp)
			}
		})
	}
}
