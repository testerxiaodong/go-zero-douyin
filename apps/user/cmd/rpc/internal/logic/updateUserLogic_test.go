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
	"gorm.io/gen"
	"gorm.io/gorm"
	"testing"
)

func TestUpdateUserLogic_UpdateUser(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockUserDo := mock.NewMockUserDo(ctl)

	serviceContext := &svc.ServiceContext{UserDo: mockUserDo}

	updateUserLogic := logic.NewUpdateUserLogic(context.Background(), serviceContext)

	// 查询用户信息失败mock
	mockUserDo.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Return(nil, errors.New("search database error"))

	// 用户不存在mock
	mockUserDo.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	// 更新用户失败mock
	mockUserDo.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Return(&model.User{Username: "test", Password: "test", ID: 1}, nil)
	mockUserDo.EXPECT().UpdateUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, errors.New("update database error"))

	// 更新用户成功
	mockUserDo.EXPECT().GetUserById(gomock.Any(), gomock.Any()).Return(&model.User{Username: "test", Password: "test", ID: 1}, nil)
	mockUserDo.EXPECT().UpdateUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(gen.ResultInfo{RowsAffected: 1, Error: nil}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.UpdateUserReq
	}{
		{
			name: "login_with_empty_param",
			req:  nil,
		},
		{
			name: "login_with_empty_id",
			req:  &pb.UpdateUserReq{Password: "test", Username: "test"},
		},
		{
			name: "login_with_search_database_error",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
		},
		{
			name: "login_with_no_database_record",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
		},
		{
			name: "login_with_update_database_error",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
		},
		{
			name: "update_success",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			infoResp, err := updateUserLogic.UpdateUser(tc.req)
			if err == nil {
				assert.NotNil(t, infoResp)
			} else {
				assert.Empty(t, infoResp)
			}
		})
	}
}
