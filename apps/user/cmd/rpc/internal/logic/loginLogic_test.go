package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/user/cmd/rpc/internal/config"
	"go-zero-douyin/apps/user/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/user/cmd/rpc/internal/model"
	"go-zero-douyin/apps/user/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/user/cmd/rpc/mock"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"testing"
)

func TestLoginLogic_Login(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserModel := mock.NewMockuserModel(ctl)

	serviceContext := &svc.ServiceContext{UserModel: mockUserModel, Config: config.Config{JwtAuth: struct {
		AccessSecret string
		AccessExpire int64
	}{AccessSecret: "test", AccessExpire: 600}}}

	loginLogic := logic.NewLoginLogic(context.Background(), serviceContext)
	dbError := errors.New("get user by username failed")
	// 查询失败mock
	mockUserModel.EXPECT().FindOneByUsernameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, dbError)
	// 用户不存在mock
	mockUserModel.EXPECT().FindOneByUsernameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, model.ErrNotFound)
	// 密码不匹配mock
	mockUserModel.EXPECT().FindOneByUsernameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.User{Password: "test"}, nil)
	// 成功的mock
	mockUserModel.EXPECT().FindOneByUsernameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.User{Password: utils.Md5ByString("test")}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.LoginReq
		err  error
	}{
		{
			name: "login_with_empty_param",
			req:  nil,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "empty param for user login"),
		},
		{
			name: "login_with_database_search_error",
			req:  &pb.LoginReq{Username: "test", Password: "test"},
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"根据用户名查询用户失败, err: %v, username: %s", dbError, "test"),
		},
		{
			name: "login_with_no_database_record",
			req:  &pb.LoginReq{Username: "test", Password: "test"},
			err:  errors.Wrapf(logic.ErrUserNotFound, "username: %s", "test"),
		},
		{
			name: "login_with_no_password_error",
			req:  &pb.LoginReq{Username: "test", Password: "test"},
			err:  errors.Wrap(logic.ErrUsernamePwdError, "密码不匹配"),
		},
		{
			name: "login_success",
			req:  &pb.LoginReq{Username: "test", Password: "test"},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := loginLogic.Login(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}

}
