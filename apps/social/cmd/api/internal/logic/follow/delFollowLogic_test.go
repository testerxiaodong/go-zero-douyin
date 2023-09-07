package follow_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/social/cmd/api/internal/logic/follow"
	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"
	socialMock "go-zero-douyin/apps/social/cmd/rpc/mock"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"go-zero-douyin/mock"
	"testing"
)

func TestDelFollowLogic_DelFollow(t *testing.T) {
	ctl := gomock.NewController(t)

	mockValidator := mock.NewMockValidator(ctl)

	mockSocialRpc := socialMock.NewMockSocial(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, SocialRpc: mockSocialRpc}

	delFollowLogic := follow.NewDelFollowLogic(context.Background(), serviceContext)

	// 参数校验失败mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// 参数校验成功，SocialRpc.UnfollowUser失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	socialRpcError := errors.New("SocialRpc.UnfollowUser error")
	mockSocialRpc.EXPECT().UnfollowUser(gomock.Any(), gomock.Any()).Return(nil, socialRpcError)

	// 参数校验成功，SocialRpc.UnfollowUser成功的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockSocialRpc.EXPECT().UnfollowUser(gomock.Any(), gomock.Any()).Return(&socialPb.UnfollowUserResp{}, nil)

	testCases := []struct {
		name string
		req  *types.UserUnfollowReq
		err  error
	}{
		{
			name: "user_unfollow_with_validate_error",
			req:  &types.UserUnfollowReq{UserId: 1},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "user_unfollow_with_social_rpc_error",
			req:  &types.UserUnfollowReq{UserId: 1},
			err:  errors.Wrapf(socialRpcError, "req: %v", &types.UserUnfollowReq{UserId: 1}),
		},
		{
			name: "user_unfollow_success",
			req:  &types.UserUnfollowReq{UserId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := delFollowLogic.DelFollow(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
