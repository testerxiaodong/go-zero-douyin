package like_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/social/cmd/api/internal/logic/like"
	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"
	socialMock "go-zero-douyin/apps/social/cmd/rpc/mock"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"go-zero-douyin/mock"
	"testing"
)

func TestUserLikeVideoIdListLogic_UserLikeVideoIdList(t *testing.T) {
	ctl := gomock.NewController(t)

	mockValidator := mock.NewMockValidator(ctl)

	mockSocialRpc := socialMock.NewMockSocial(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, SocialRpc: mockSocialRpc}

	userLikeVideoIdListLogic := like.NewUserLikeVideoIdListLogic(context.Background(), serviceContext)

	// 参数校验失败mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// 参数校验成功，但SocialRpc.GetUserLikeVideoIdList失败mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	socialRpcError := errors.New("SocialRpc.GetUserLikeVideoIdList error")
	mockSocialRpc.EXPECT().GetUserLikeVideoIdList(gomock.Any(), gomock.Any()).Return(&socialPb.GetUserLikeVideoIdListResp{}, socialRpcError)

	// 参数校验成功，且SocialRpc.GetUserLikeVideoIdList调用成功的mock
	expectedResp := &socialPb.GetUserLikeVideoIdListResp{
		VideoIdList: []int64{1, 2, 3}}
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockSocialRpc.EXPECT().GetUserLikeVideoIdList(gomock.Any(), gomock.Any()).Return(expectedResp, nil)

	testCases := []struct {
		name string
		req  *types.GetUserLikeVideoIdListReq
		err  error
	}{
		{
			name: "user_like_video_id_list_with_validate_error",
			req:  &types.GetUserLikeVideoIdListReq{UserId: 1},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "user_like_video_id_list_with_social_rpc_error",
			req:  &types.GetUserLikeVideoIdListReq{UserId: 1},
			err:  errors.Wrapf(socialRpcError, "req: %v", &types.GetUserLikeVideoIdListReq{UserId: 1}),
		},
		{
			name: "user_like_video_id_list_success",
			req:  &types.GetUserLikeVideoIdListReq{UserId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := userLikeVideoIdListLogic.UserLikeVideoIdList(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
				assert.Equal(t, len(resp.VideoIdList), 3)
			}
		})
	}
}
