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

func TestGetUserFollowCountLogic_GetUserFollowCount(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockFollowDo := mock.NewMockfollowCountModel(ctl)
	serviceContext := &svc.ServiceContext{FollowCountModel: mockFollowDo}
	getUserFollowCountLogic := logic.NewGetUserFollowCountLogic(context.Background(), serviceContext)

	// 查询数据库失败的mock
	dbError := errors.New("FollowDo GetUserFollowCount error")
	mockFollowDo.EXPECT().FindOneByUserIdIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbError)

	// 查询数据库成功，但没有数据的的mock
	mockFollowDo.EXPECT().FindOneByUserIdIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, model.ErrNotFound)

	// 查询数据库成功，有数据的的mock
	mockFollowDo.EXPECT().FindOneByUserIdIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.FollowCount{FollowCount: 1}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.GetUserFollowCountReq
		err  error
	}{
		{
			name: "get_user_follow_count_with_empty_param",
			req:  nil,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get user follow count with empty param"),
		},
		{
			name: "get_user_follow_count_with_empty_video_id",
			req:  &pb.GetUserFollowCountReq{UserId: 0},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get user follow count with empty user_id"),
		},
		{
			name: "get_user_follow_count_with_database_error",
			req:  &pb.GetUserFollowCountReq{UserId: 10},
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"查询用户关注数粉丝数失败, err: %v, user_id: %d", dbError, 10),
		},
		{
			name: "get_user_follow_count_with_no_database_record",
			req:  &pb.GetUserFollowCountReq{UserId: 10},
			err:  nil,
		},
		{
			name: "get_user_follow_count_with_database_record",
			req:  &pb.GetUserFollowCountReq{UserId: 10},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getUserFollowCountLogic.GetUserFollowCount(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
