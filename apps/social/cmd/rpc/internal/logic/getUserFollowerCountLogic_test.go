package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/social/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/mock"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"
	"testing"
)

func TestGetUserFollowerCountLogic_GetUserFollowerCount(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockFollowDo := mock.NewMockfollowCountModel(ctl)
	serviceContext := &svc.ServiceContext{FollowCountModel: mockFollowDo}
	getUserFollowerCountLogic := logic.NewGetUserFollowerCountLogic(context.Background(), serviceContext)

	// 查询数据库失败的mock
	dbError := errors.New("FollowDo GetUserFollowCount error")
	mockFollowDo.EXPECT().FindOneByUserIdIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbError)

	// 查询数据库成功，但没有数据的mock
	mockFollowDo.EXPECT().FindOneByUserIdIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, model.ErrNotFound)

	// 查询数据库成功，有数据的mock
	mockFollowDo.EXPECT().FindOneByUserIdIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.FollowCount{FollowerCount: 1}, model.ErrNotFound)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.GetUserFollowerCountReq
		err  error
	}{
		{
			name: "get_user_follower_count_with_empty_param",
			req:  nil,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follower count with empty param"),
		},
		{
			name: "get_user_follower_count_with_empty_video_id",
			req:  &pb.GetUserFollowerCountReq{UserId: 0},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follower count with empty user_id"),
		},
		{
			name: "get_user_follower_count_with_database_error",
			req:  &pb.GetUserFollowerCountReq{UserId: 10},
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"查询用户粉丝数关注数失败, err: %v, user_id: %d", dbError, 10),
		},
		{
			name: "get_user_follower_count_with_no_database_record",
			req:  &pb.GetUserFollowerCountReq{UserId: 10},
			err:  nil,
		},
		{
			name: "get_user_follower_count_with_database_record",
			req:  &pb.GetUserFollowerCountReq{UserId: 10},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getUserFollowerCountLogic.GetUserFollowerCount(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
