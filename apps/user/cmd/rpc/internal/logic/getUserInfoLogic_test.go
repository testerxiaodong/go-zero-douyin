package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/user/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/user/cmd/rpc/internal/model"
	"go-zero-douyin/apps/user/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/user/cmd/rpc/mock"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	"gorm.io/gorm"
	"testing"
)

func TestGetUserInfoLogic_GetUserInfo(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockUserDo := mock.NewMockUserDo(ctl)

	serviceContext := &svc.ServiceContext{UserDo: mockUserDo}

	getUserInfoLogic := logic.NewGetUserInfoLogic(context.Background(), serviceContext)

	mockUserDo.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Return(&model.User{ID: 2, Username: "test"}, nil)
	mockUserDo.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Return(nil, errors.New("database search error"))
	mockUserDo.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.GetUserInfoReq
	}{
		{
			name: "get_user_info_success",
			req:  &pb.GetUserInfoReq{Id: 2},
		},
		{
			name: "get_user_info_with_empty_param",
			req:  nil,
		},
		{
			name: "get_user_info_with_database_search_error",
			req:  &pb.GetUserInfoReq{Id: 2},
		},
		{
			name: "get_user_info_with_no_database_record",
			req:  &pb.GetUserInfoReq{Id: 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			infoResp, err := getUserInfoLogic.GetUserInfo(tc.req)
			if err == nil {
				assert.NotEmpty(t, infoResp)
			} else {
				assert.Empty(t, infoResp)
			}
		})
	}
}
