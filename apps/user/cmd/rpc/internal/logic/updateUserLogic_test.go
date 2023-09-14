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
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"gorm.io/gen"
	"testing"
)

func TestUpdateUserLogic_UpdateUser(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockUserDo := mock.NewMockUserDo(ctl)

	mockSender := globalMock.NewMockSender(ctl)

	serviceContext := &svc.ServiceContext{UserDo: mockUserDo, Rabbit: mockSender}

	updateUserLogic := logic.NewUpdateUserLogic(context.Background(), serviceContext)

	// 查询用户信息失败mock
	dbError := errors.New("search database error")
	mockUserDo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Return(nil, dbError)

	// 用户名已存在mock
	mockUserDo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Return(&model.User{ID: 2}, nil)

	// 更新用户失败mock
	dbUpdateError := errors.New("update database error")
	mockUserDo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Return(&model.User{Username: "test", Password: "test", ID: 1}, nil)
	mockUserDo.EXPECT().UpdateUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, dbUpdateError)

	// 更新用户成功，但发送消息失败的mock
	senderError := errors.New("sender error")
	mockUserDo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Return(&model.User{Username: "test", Password: "test", ID: 1}, nil)
	mockUserDo.EXPECT().UpdateUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(gen.ResultInfo{RowsAffected: 1, Error: nil}, nil)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 更新用户成功的mock
	mockUserDo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Return(&model.User{Username: "test", Password: "test", ID: 1}, nil)
	mockUserDo.EXPECT().UpdateUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(gen.ResultInfo{RowsAffected: 1, Error: nil}, nil)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.UpdateUserReq
		err  error
	}{
		{
			name: "update_user_with_empty_param",
			req:  nil,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Update user empty param"),
		},
		{
			name: "update_user_with_empty_id",
			req:  &pb.UpdateUserReq{Password: "test", Username: "test"},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Updtae user empty user id"),
		},
		{
			name: "update_user_with_search_database_error",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "Find user by id failed, user_id: %d, err: %v", 1, dbError),
		},
		{
			name: "update_user_with_username_exist",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
			err:  errors.Wrapf(logic.ErrUserAlreadyRegister, "username: %s", "test"),
		},
		{
			name: "update_user_with_update_database_error",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR), "update mysql user info failed, err: %v", dbUpdateError),
		},
		{
			name: "update_user_with_sender_error",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.RPC_UPDATE_ERR), "发布更新es用户文档信息失败, err: %v", senderError),
		},
		{
			name: "update_success",
			req:  &pb.UpdateUserReq{Username: "test", Password: "test", Id: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := updateUserLogic.UpdateUser(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
