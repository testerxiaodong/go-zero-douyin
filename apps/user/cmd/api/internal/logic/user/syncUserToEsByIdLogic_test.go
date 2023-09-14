package user_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	searchMock "go-zero-douyin/apps/search/cmd/rpc/mock"
	searchPb "go-zero-douyin/apps/search/cmd/rpc/pb"
	socialMock "go-zero-douyin/apps/social/cmd/rpc/mock"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/apps/user/cmd/api/internal/logic/user"
	"go-zero-douyin/apps/user/cmd/api/internal/svc"
	"go-zero-douyin/apps/user/cmd/api/internal/types"
	"go-zero-douyin/apps/user/cmd/rpc/mock"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	gloablMock "go-zero-douyin/mock"
	"testing"
)

func TestSyncUserToEsByIdLogic_SyncUserToEsById(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockValidator := gloablMock.NewMockValidator(ctl)
	mockUserRpc := mock.NewMockUser(ctl)
	mockSocialRpc := socialMock.NewMockSocial(ctl)
	mockSearchRpc := searchMock.NewMockSearch(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, UserRpc: mockUserRpc,
		SearchRpc: mockSearchRpc, SocialRpc: mockSocialRpc}

	syncUserToEsByIdLogic := user.NewSyncUserToEsByIdLogic(context.Background(), serviceContext)

	// 参数校验失败的mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// userRpc失败的mock
	userRpcError := errors.New("userRpc.GetUserInfo error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpc.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(nil, userRpcError)

	// socialRpc.GetUserFollowerCount失败的mock
	followerCountError := errors.New("socialRpc.GetUserFollowerCount error")
	userInfo := &pb.UserInfo{
		Id:       utils.NewRandomInt64(1, 10),
		Username: utils.NewRandomString(10),
	}
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpc.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&pb.GetUserInfoResp{User: userInfo}, nil)
	mockSocialRpc.EXPECT().GetUserFollowerCount(gomock.Any(), gomock.Any()).Return(nil, followerCountError)

	// socialRpc.GetUserFollowCount失败的mock
	followCountError := errors.New("socialRpc.GetUserFollowCount error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpc.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&pb.GetUserInfoResp{User: userInfo}, nil)
	mockSocialRpc.EXPECT().GetUserFollowerCount(gomock.Any(), gomock.Any()).Return(&socialPb.GetUserFollowerCountResp{FollowerCount: 1}, nil)
	mockSocialRpc.EXPECT().GetUserFollowCount(gomock.Any(), gomock.Any()).Return(nil, followCountError)

	// SearchRpc.SyncUserInfo失败的mock
	searchRpcError := errors.New("SearchRpc.SyncUserInfo error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpc.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&pb.GetUserInfoResp{User: userInfo}, nil)
	mockSocialRpc.EXPECT().GetUserFollowerCount(gomock.Any(), gomock.Any()).Return(&socialPb.GetUserFollowerCountResp{FollowerCount: 1}, nil)
	mockSocialRpc.EXPECT().GetUserFollowCount(gomock.Any(), gomock.Any()).Return(&socialPb.GetUserFollowCountResp{FollowCount: 1}, nil)
	mockSearchRpc.EXPECT().SyncUserInfo(gomock.Any(), gomock.Any()).Return(nil, searchRpcError)

	// 成功的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockUserRpc.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).Return(&pb.GetUserInfoResp{User: userInfo}, nil)
	mockSocialRpc.EXPECT().GetUserFollowerCount(gomock.Any(), gomock.Any()).Return(&socialPb.GetUserFollowerCountResp{FollowerCount: 1}, nil)
	mockSocialRpc.EXPECT().GetUserFollowCount(gomock.Any(), gomock.Any()).Return(&socialPb.GetUserFollowCountResp{FollowCount: 1}, nil)
	mockSearchRpc.EXPECT().SyncUserInfo(gomock.Any(), gomock.Any()).Return(&searchPb.SyncUserInfoResp{}, nil)

	// 表格启动测试
	req := &types.SyncUserToEsByIdReq{UserId: utils.NewRandomInt64(1, 10)}
	testCases := []struct {
		name string
		req  *types.SyncUserToEsByIdReq
		err  error
	}{
		{
			name: "sync_user_to_es_with_validate_error",
			req:  req,
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "sync_user_to_es_with_user_rpc_error",
			req:  req,
			err:  errors.Wrapf(userRpcError, "req: %v", req),
		},
		{
			name: "sync_user_to_es_with_follower_count_error",
			req:  req,
			err:  errors.Wrapf(followerCountError, "req: %v", req),
		},
		{
			name: "sync_user_to_es_with_follow_count_error",
			req:  req,
			err:  errors.Wrapf(followCountError, "req: %v", req),
		},
		{
			name: "sync_user_to_es_with_search_rpc_error",
			req:  req,
			err:  errors.Wrapf(searchRpcError, "req: %v", req),
		},
		{
			name: "sync_user_to_es_success",
			req:  req,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := syncUserToEsByIdLogic.SyncUserToEsById(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
