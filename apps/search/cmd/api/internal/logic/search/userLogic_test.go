package search_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/search/cmd/api/internal/logic/search"
	"go-zero-douyin/apps/search/cmd/api/internal/svc"
	"go-zero-douyin/apps/search/cmd/api/internal/types"
	searchMock "go-zero-douyin/apps/search/cmd/rpc/mock"
	"go-zero-douyin/apps/search/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestUserLogic_User(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockValidator := globalMock.NewMockValidator(ctl)
	mockSearchRpc := searchMock.NewMockSearch(ctl)
	serviceContext := &svc.ServiceContext{Validator: mockValidator, SearchRpc: mockSearchRpc}
	userLogic := search.NewUserLogic(context.Background(), serviceContext)

	// 参数校验失败的mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// searchRpc.SearchUser失败的mock
	esError := errors.New("es error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockSearchRpc.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return(nil, esError)

	// 成功的mock
	resp := &pb.SearchUserResp{Total: 1, Users: []*pb.User{&pb.User{
		Id:            utils.NewRandomInt64(1, 10),
		Username:      utils.NewRandomString(10),
		FollowerCount: utils.NewRandomInt64(1, 10),
		FollowCount:   utils.NewRandomInt64(1, 10),
	}}}
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockSearchRpc.EXPECT().SearchUser(gomock.Any(), gomock.Any()).Return(resp, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *types.SearchUserReq
		err  error
	}{
		{
			name: "search_user_with_validate_error",
			req:  &types.SearchUserReq{Keyword: "test", Sort: 0},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "search_user_with_es_error",
			req:  &types.SearchUserReq{Keyword: "test", Sort: 0},
			err:  errors.Wrapf(esError, "req: %v", &types.SearchUserReq{Keyword: "test", Sort: 0}),
		},
		{
			name: "search_user_success",
			req:  &types.SearchUserReq{Keyword: "test", Sort: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := userLogic.User(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
