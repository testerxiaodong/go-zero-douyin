package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/syncx"
	"go-zero-douyin/apps/social/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/mock"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestGetUserFollowerCountLogic_GetUserFollowerCount(t *testing.T) {
	ctl := gomock.NewController(t)

	mockFollowDo := mock.NewMockFollowDo(ctl)

	mockRedis := globalMock.NewMockRedisCache(ctl)

	utils.IgnoreGo()
	defer utils.RecoverGo()

	serviceContext := &svc.ServiceContext{FollowDo: mockFollowDo, Redis: mockRedis, SingleFlight: syncx.NewSingleFlight()}

	getUserFollowerCountLogic := logic.NewGetUserFollowerCountLogic(context.Background(), serviceContext)

	// redis中有数据mock
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
	mockRedis.EXPECT().Scard(gomock.Any(), gomock.Any()).Return(int64(10), nil)
	mockRedis.EXPECT().Expire(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// redis中没有数据，查询数据库失败的mock
	dbError := errors.New("FollowDo GetUserFollowCount error")
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
	mockFollowDo.EXPECT().GetUserFollowerCount(gomock.Any(), gomock.Any()).Return(int64(0), dbError)

	// redis中没有数据，查询数据库成功的mock
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
	mockFollowDo.EXPECT().GetUserFollowerCount(gomock.Any(), gomock.Any()).Return(int64(10), nil)

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
			name: "get_user_follower_count_with_redis",
			req:  &pb.GetUserFollowerCountReq{UserId: 10},
			err:  nil,
		},
		{
			name: "get_user_follower_count_with_database_error",
			req:  &pb.GetUserFollowerCountReq{UserId: 10},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "get user follower count from mysql failed, err: %v, user_id: %d", dbError, 10),
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
