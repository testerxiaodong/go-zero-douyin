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
	gloabelMock "go-zero-douyin/mock"
	"gorm.io/gen"
	"gorm.io/gorm"
	"testing"
)

func TestUnfollowUserLogic_UnfollowUser(t *testing.T) {
	ctl := gomock.NewController(t)

	mockFollowDo := mock.NewMockFollowDo(ctl)

	mockSender := gloabelMock.NewMockSender(ctl)

	mockRedis := gloabelMock.NewMockRedisCache(ctl)

	serviceContext := &svc.ServiceContext{FollowDo: mockFollowDo, Redis: mockRedis, Rabbit: mockSender}

	unfollowUserLogic := logic.NewUnfollowUserLogic(context.Background(), serviceContext)

	// 查询数据库失败mock
	dbSearchError := errors.New("FollowDo.GetFollowByFollowerIdAndUserId error")
	mockFollowDo.EXPECT().GetFollowByFollowerIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 查询数据库成功，数据不存在的mock
	mockFollowDo.EXPECT().GetFollowByFollowerIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	// 查询数据库成功，数据存在，但删除数据库失败mock
	dbDeleteError := errors.New("FollowDo.DeleteFollow error")
	mockFollowDo.EXPECT().GetFollowByFollowerIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Follow{}, nil)
	mockFollowDo.EXPECT().DeleteFollow(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, dbDeleteError)

	// 查询数据库成功，数据存在，删除数据成功，删除关注缓存失败，发布关注消息失败
	redisError := errors.New("redis delete error")
	senderError := errors.New("rabbitmq sender error")
	mockFollowDo.EXPECT().GetFollowByFollowerIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Follow{}, nil)
	mockFollowDo.EXPECT().DeleteFollow(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 查询数据库成功，数据存在，删除数据成功，删除关注缓存失败，发布关注消息成功，删除粉丝缓存失败，发送粉丝消息失败的mock
	mockFollowDo.EXPECT().GetFollowByFollowerIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Follow{}, nil)
	mockFollowDo.EXPECT().DeleteFollow(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 查询数据库成功，数据存在，删除数据成功，删除关注缓存失败，发布关注消息成功，删除粉丝缓存失败，发送粉丝消息成功的mock
	mockFollowDo.EXPECT().GetFollowByFollowerIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Follow{}, nil)
	mockFollowDo.EXPECT().DeleteFollow(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 查询数据库成功，数据存在，删除数据成功，删除关注缓存失败，发布关注消息成功，删除粉丝缓存成功的mock
	mockFollowDo.EXPECT().GetFollowByFollowerIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Follow{}, nil)
	mockFollowDo.EXPECT().DeleteFollow(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)

	// 查询数据库成功，数据存在，删除数据成功，删除关注缓存成功，删除粉丝缓存成功的mock
	mockFollowDo.EXPECT().GetFollowByFollowerIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Follow{}, nil)
	mockFollowDo.EXPECT().DeleteFollow(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.UnfollowUserReq
		err  error
	}{
		{
			name: "unfollow_user_with_empty_param",
			req:  nil,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unfollow user with empty param"),
		},
		{
			name: "unfollow_user_with_empty_filed",
			req:  &pb.UnfollowUserReq{FollowerId: 0, UserId: 0},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unfollow user with empty follower_id or user_id"),
		},
		{
			name: "unfollow_user_self_error",
			req:  &pb.UnfollowUserReq{FollowerId: 1, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("不能对自己操作"), "req: %v", &pb.FollowUserReq{FollowerId: 1, UserId: 1}),
		},
		{
			name: "unfollow_user_with_database_search_error",
			req:  &pb.UnfollowUserReq{FollowerId: 2, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "find follow record by follower_id and user_id failed, err: %v, follower_id: %d, user_id: %d", dbSearchError, 2, 1),
		},
		{
			name: "unfollow_user_with_database_record_exist",
			req:  &pb.UnfollowUserReq{FollowerId: 2, UserId: 1},
			err:  nil,
		},
		{
			name: "unfollow_user_with_database_delete_error",
			req:  &pb.UnfollowUserReq{FollowerId: 2, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "delete follow record failed, err: %v, follower_id: %d, user_id: %d", dbDeleteError, 2, 1),
		},
		{
			name: "unfollow_user_with_follow_rabbit_error",
			req:  &pb.UnfollowUserReq{FollowerId: 2, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("发布userFollowUserMessage失败"), "err: %v", senderError),
		},
		{
			name: "unfollow_user_with_follower_rabbit_error",
			req:  &pb.UnfollowUserReq{FollowerId: 2, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("发布userFollowedByUserMessage失败"), "err: %v", senderError),
		},
		{
			name: "unfollow_user_with_redis_error1",
			req:  &pb.UnfollowUserReq{FollowerId: 2, UserId: 1},
			err:  nil,
		},
		{
			name: "unfollow_user_with_redis_error2",
			req:  &pb.UnfollowUserReq{FollowerId: 2, UserId: 1},
			err:  nil,
		},
		{
			name: "unfollow_user_success",
			req:  &pb.UnfollowUserReq{FollowerId: 2, UserId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := unfollowUserLogic.UnfollowUser(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
