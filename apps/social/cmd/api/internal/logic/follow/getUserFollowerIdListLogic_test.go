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

func TestGetUserFollowerIdListLogic_GetUserFollowerIdList(t *testing.T) {
	ctl := gomock.NewController(t)

	mockValidator := mock.NewMockValidator(ctl)

	mockSocialRpc := socialMock.NewMockSocial(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, SocialRpc: mockSocialRpc}

	getUserFollowerIdListLogic := follow.NewGetUserFollowerIdListLogic(context.Background(), serviceContext)

	// 参数校验失败mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// 参数校验成功，SocialRpc.GetUserFollowerIdList失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	socialRpcError := errors.New("SocialRpc.GetUserFollowerIdList error")
	mockSocialRpc.EXPECT().GetUserFollowerIdList(gomock.Any(), gomock.Any()).Return(nil, socialRpcError)

	// 参数校验成功，SocialRpc.GetUserFollowerIdList成功的mock
	expectedResp := []int64{utils.NewRandomInt64(1, 100), utils.NewRandomInt64(1, 100)}
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockSocialRpc.EXPECT().GetUserFollowerIdList(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetUserFollowerIdListResp{UserIdList: expectedResp}, nil)

	testCases := []struct {
		name string
		req  *types.GetUserFollowerIdListReq
		err  error
	}{
		{
			name: "get_user_follower_id_list_with_validate_error",
			req:  &types.GetUserFollowerIdListReq{UserId: 1},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "get_user_follower_id_list_with_social_rpc_error",
			req:  &types.GetUserFollowerIdListReq{UserId: 1},
			err:  errors.Wrapf(socialRpcError, "req: %v", &types.GetUserFollowerIdListReq{UserId: 1}),
		},
		{
			name: "get_user_follower_id_list_success",
			req:  &types.GetUserFollowerIdListReq{UserId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := getUserFollowerIdListLogic.GetUserFollowerIdList(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
				assert.Equal(t, resp.UserIdList, expectedResp)
			}
		})
	}
}
