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
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestGetUserFollowerIdListLogic_GetUserFollowerIdList(t *testing.T) {
	ctl := gomock.NewController(t)

	mockFollowDo := mock.NewMockFollowDo(ctl)

	mockRedis := globalMock.NewMockRedisCache(ctl)

	serviceContext := &svc.ServiceContext{FollowDo: mockFollowDo, Redis: mockRedis, SingleFlight: syncx.NewSingleFlight()}

	getUserFollowerIdListLogic := logic.NewGetUserFollowerIdListLogic(context.Background(), serviceContext)

	// redis中有数据mock
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
	mockRedis.EXPECT().Smembers(gomock.Any(), gomock.Any()).Return([]string{"1", "2", "3"}, nil)
	mockRedis.EXPECT().Expire(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// redis中没有数据，查询数据库失败的mock
	dbError := errors.New("FollowDo GetUserFollowerIdList error")
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
	mockFollowDo.EXPECT().GetUserFollowerIdList(gomock.Any(), gomock.Any()).Return([]int64{}, dbError)

	// redis中没有数据，查询数据库成功的mock
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
	mockFollowDo.EXPECT().GetUserFollowerIdList(gomock.Any(), gomock.Any()).Return([]int64{1, 2, 3}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.GetUserFollowerIdListReq
		err  error
	}{
		{
			name: "get_user_follower_id_list_with_empty_param",
			req:  nil,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follower id list with empty param"),
		},
		{
			name: "get_user_follower_id_list_with_empty_video_id",
			req:  &pb.GetUserFollowerIdListReq{UserId: 0},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follower id list with empty user_id"),
		},
		{
			name: "get_user_follower_id_list_with_redis",
			req:  &pb.GetUserFollowerIdListReq{UserId: 10},
			err:  nil,
		},
		{
			name: "get_user_follower_id_list_with_database_error",
			req:  &pb.GetUserFollowerIdListReq{UserId: 10},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "get user follower id list from mysql failed, err: %v", dbError),
		},
		{
			name: "get_user_follower_id_list_with_database_record",
			req:  &pb.GetUserFollowerIdListReq{UserId: 10},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getUserFollowerIdListLogic.GetUserFollowerIdList(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
