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

func TestVideoUnlikeLogic_VideoUnlike(t *testing.T) {
	ctl := gomock.NewController(t)

	mockLikeDo := mock.NewMockLikeDo(ctl)

	mockSender := gloabelMock.NewMockSender(ctl)

	mockRedis := gloabelMock.NewMockRedisCache(ctl)

	serviceContext := &svc.ServiceContext{LikeDo: mockLikeDo, Redis: mockRedis, Rabbit: mockSender}

	videoUnlikeLogic := logic.NewVideoUnlikeLogic(context.Background(), serviceContext)

	// 查询数据库失败mock
	dbSearchError := errors.New("LikeDo.GetLikeByVideoIdAndUserId error")
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 查询数据库成功，数据不存在的mock
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	// 查询数据库成功，数据存在，但删除数据库失败mock
	dbDeleteError := errors.New("LikeDo.DeleteLike error")
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Like{}, nil)
	mockLikeDo.EXPECT().DeleteLike(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, dbDeleteError)

	// 查询数据库成功，数据存在，删除数据成功，删除关注缓存失败，发布关注消息失败
	redisError := errors.New("redis delete error")
	senderError := errors.New("rabbitmq sender error")
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Like{}, nil)
	mockLikeDo.EXPECT().DeleteLike(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 查询数据库成功，数据存在，删除数据成功，删除关注缓存失败，发布关注消息成功，删除粉丝缓存失败，发送粉丝消息失败的mock
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Like{}, nil)
	mockLikeDo.EXPECT().DeleteLike(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 查询数据库成功，数据存在，删除数据成功，删除关注缓存失败，发布关注消息成功，删除粉丝缓存失败，发送粉丝消息成功的mock
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Like{}, nil)
	mockLikeDo.EXPECT().DeleteLike(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 查询数据库成功，数据存在，删除数据成功，删除关注缓存失败，发布关注消息成功，删除粉丝缓存成功的mock
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Like{}, nil)
	mockLikeDo.EXPECT().DeleteLike(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)

	// 查询数据库成功，数据存在，删除数据成功，删除关注缓存成功，删除粉丝缓存成功的mock
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Like{}, nil)
	mockLikeDo.EXPECT().DeleteLike(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.VideoUnlikeReq
		err  error
	}{
		{
			name: "video_unlike_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unlike video with empty param"),
		},
		{
			name: "video_unlike_with_empty_filed",
			req:  &pb.VideoUnlikeReq{VideoId: 0, UserId: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unlike video with empty video_id or user_id"),
		},
		{
			name: "video_unlike_with_database_search_error",
			req:  &pb.VideoUnlikeReq{VideoId: 2, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "find video is liked by user failed, err: %v", dbSearchError),
		},
		{
			name: "video_unlike_with_database_record_exist",
			req:  &pb.VideoUnlikeReq{VideoId: 2, UserId: 1},
			err:  nil,
		},
		{
			name: "video_unlike_with_database_delete_error",
			req:  &pb.VideoUnlikeReq{VideoId: 2, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.RPC_DELETE_ERR), "delete user like video record failed, err: %v", dbDeleteError),
		},
		{
			name: "video_unlike_with_follow_rabbit_error",
			req:  &pb.VideoUnlikeReq{VideoId: 2, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("publish user like video message failed"), "video_id: %d", 2),
		},
		{
			name: "video_unlike_with_follower_rabbit_error",
			req:  &pb.VideoUnlikeReq{VideoId: 2, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("publish video liked by user message failed"), "user_id: %d", 1),
		},
		{
			name: "video_unlike_with_redis_error1",
			req:  &pb.VideoUnlikeReq{VideoId: 2, UserId: 1},
			err:  nil,
		},
		{
			name: "video_unlike_with_redis_error2",
			req:  &pb.VideoUnlikeReq{VideoId: 2, UserId: 1},
			err:  nil,
		},
		{
			name: "video_unlike_success",
			req:  &pb.VideoUnlikeReq{VideoId: 2, UserId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := videoUnlikeLogic.VideoUnlike(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
