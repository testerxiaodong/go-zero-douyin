package user_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	socialMock "go-zero-douyin/apps/social/cmd/rpc/mock"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
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

func TestDetailLogic_Detail(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	// user rpc client的mock
	mockUserRpcClient := mock.NewMockUser(ctl)

	mockSocialRpc := socialMock.NewMockSocial(ctl)

	// validator mock
	mockValidator := globalMock.NewMockValidator(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, UserRpc: mockUserRpcClient, SocialRpc: mockSocialRpc}

	detailLogic := user.NewDetailLogic(context.Background(), serviceContext)

	// 期望的结果
	expectedUserInfo := &pb.UserInfo{
		Id:       utils.NewRandomInt64(1, 10),
		Username: utils.NewRandomString(10),
	}

	// 参数校验失败的mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// 参数校验成功，但userRpc.GetUserInfo失败的mock
	userRpcError := errors.New("userRpc.GetUserInfo error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpcClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).
		Return(nil, userRpcError)

	// SocialRpc.GetUserFollowerCount失败的mock
	followerCountError := errors.New("SocialRpc.GetUserFollowerCount error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpcClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).
		Return(&pb.GetUserInfoResp{User: expectedUserInfo}, nil)
	mockSocialRpc.EXPECT().GetUserFollowerCount(gomock.Any(), gomock.Any()).Return(nil, followerCountError)

	// SocialRpc.GetUserFollowCount失败的mock
	followCountError := errors.New("SocialRpc.GetUserFollowCount error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpcClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).
		Return(&pb.GetUserInfoResp{User: expectedUserInfo}, nil)
	mockSocialRpc.EXPECT().GetUserFollowerCount(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetUserFollowerCountResp{FollowerCount: 1}, nil)
	mockSocialRpc.EXPECT().GetUserFollowCount(gomock.Any(), gomock.Any()).
		Return(nil, followCountError)

	// 参数校验成功，且userRpc.GetUserInfo成功
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpcClient.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).
		Return(&pb.GetUserInfoResp{User: expectedUserInfo}, nil)
	mockSocialRpc.EXPECT().GetUserFollowerCount(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetUserFollowerCountResp{FollowerCount: 1}, nil)
	mockSocialRpc.EXPECT().GetUserFollowCount(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetUserFollowCountResp{FollowCount: 1}, nil)

	// 表格驱动测试
	req := &types.UserInfoReq{Id: 1}
	testCases := []struct {
		name string
		req  *types.UserInfoReq
		err  error
	}{
		{
			name: "get_user_detail_with_validator_error",
			req:  req,
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "get_user_detail_with_user_rpc_error",
			req:  req,
			err:  errors.Wrapf(userRpcError, "req: %v", req),
		},
		{
			name: "get_user_detail_with_follower_count_error",
			req:  req,
			err:  errors.Wrapf(followerCountError, "req: %v", req),
		},
		{
			name: "get_user_detail_with_follow_count_error",
			req:  req,
			err:  errors.Wrapf(followCountError, "req: %v", req),
		},
		{
			name: "get_user_detail_success",
			req:  &types.UserInfoReq{Id: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := detailLogic.Detail(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
				assert.Equal(t, resp.User.Id, expectedUserInfo.Id)
				assert.Equal(t, resp.User.Username, expectedUserInfo.Username)
				assert.Equal(t, resp.User.FollowerCount, int64(1))
				assert.Equal(t, resp.User.FollowCount, int64(1))
			}
		})
	}
}
