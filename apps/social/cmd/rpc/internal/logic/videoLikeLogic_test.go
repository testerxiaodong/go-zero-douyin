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
	"gorm.io/gorm"
	"testing"
)

func TestVideoLikeLogic_VideoLike(t *testing.T) {
	ctl := gomock.NewController(t)

	mockLikeDo := mock.NewMockLikeDo(ctl)

	mockSender := gloabelMock.NewMockSender(ctl)

	mockRedis := gloabelMock.NewMockRedisCache(ctl)

	serviceContext := &svc.ServiceContext{LikeDo: mockLikeDo, Redis: mockRedis, Rabbit: mockSender}

	videoLikeLogic := logic.NewVideoLikeLogic(context.Background(), serviceContext)

	// 查询数据库失败mock
	dbSearchError := errors.New("LikeDo.GetLikeByVideoIdAndUserId error")
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 查询数据库成功，数据存在的mock
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Like{}, nil)

	// 查询数据库成功，数据不存在，但插入数据库失败mock
	dbInsertError := errors.New("LikeDo.InsertLike error")
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockLikeDo.EXPECT().InsertLike(gomock.Any(), gomock.Any()).Return(dbInsertError)

	// 查询数据库成功，数据不存在，插入数据成功，删除关注缓存失败，发布关注消息失败
	redisError := errors.New("redis delete error")
	senderError := errors.New("rabbitmq sender error")
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockLikeDo.EXPECT().InsertLike(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 查询数据库成功，数据不存在，插入数据成功，删除关注缓存失败，发布关注消息成功，删除粉丝缓存失败，发送粉丝消息失败的mock
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockLikeDo.EXPECT().InsertLike(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 查询数据库成功，数据不存在，插入数据成功，删除关注缓存失败，发布关注消息成功，删除粉丝缓存失败，发送粉丝消息成功的mock
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockLikeDo.EXPECT().InsertLike(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 查询数据库成功，数据不存在，插入数据成功，删除关注缓存失败，发布关注消息成功，删除粉丝缓存成功的mock
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockLikeDo.EXPECT().InsertLike(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)

	// 查询数据库成功，数据不存在，插入数据成功，删除关注缓存成功，删除粉丝缓存成功的mock
	mockLikeDo.EXPECT().GetLikeByVideoIdAndUserId(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockLikeDo.EXPECT().InsertLike(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.VideoLikeReq
		err  error
	}{
		{
			name: "video_like_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "like video with empty param"),
		},
		{
			name: "video_like_with_empty_filed",
			req:  &pb.VideoLikeReq{VideoId: 0, UserId: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "like video with empty user_id or video_id"),
		},
		{
			name: "video_like_with_database_search_error",
			req:  &pb.VideoLikeReq{VideoId: 2, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "find video is already liked by user failed, err: %v", dbSearchError),
		},
		{
			name: "video_like_with_database_record_exist",
			req:  &pb.VideoLikeReq{VideoId: 2, UserId: 1},
			err:  nil,
		},
		{
			name: "video_like_with_database_insert_error",
			req:  &pb.VideoLikeReq{VideoId: 2, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.RPC_INSERT_ERR), "insert video like failed, err: %v", dbInsertError),
		},
		{
			name: "video_like_with_follow_rabbit_error",
			req:  &pb.VideoLikeReq{VideoId: 2, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("publish user like video message failed"), "video_id: %d", 2),
		},
		{
			name: "video_like_with_follower_rabbit_error",
			req:  &pb.VideoLikeReq{VideoId: 2, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("publish video liked by user message failed"), "user_id: %d", 1),
		},
		{
			name: "video_like_with_redis_error1",
			req:  &pb.VideoLikeReq{VideoId: 2, UserId: 1},
			err:  nil,
		},
		{
			name: "video_like_with_redis_error2",
			req:  &pb.VideoLikeReq{VideoId: 2, UserId: 1},
			err:  nil,
		},
		{
			name: "video_like_success",
			req:  &pb.VideoLikeReq{VideoId: 2, UserId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := videoLikeLogic.VideoLike(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
