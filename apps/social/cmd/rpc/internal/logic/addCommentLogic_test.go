package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/social/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/mock"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"
	gloabelMock "go-zero-douyin/mock"
	"testing"
)

func TestAddCommentLogic_AddComment(t *testing.T) {
	ctl := gomock.NewController(t)

	mockCommentDo := mock.NewMockCommentDo(ctl)

	mockSender := gloabelMock.NewMockSender(ctl)

	mockRedis := gloabelMock.NewMockRedisCache(ctl)

	serviceContext := &svc.ServiceContext{CommentDo: mockCommentDo, Redis: mockRedis, Rabbit: mockSender}

	addCommentLogic := logic.NewAddCommentLogic(context.Background(), serviceContext)

	// CommentDo.InsertComment失败的mock
	insertCommentError := errors.New("CommentDo.InsertComment error")
	mockCommentDo.EXPECT().InsertComment(gomock.Any(), gomock.Any()).Return(insertCommentError)

	// 插入数据库成功，但删除缓存失败，且发布消息失败mock
	redisError := errors.New("redis delete error")
	senderError := errors.New("rabbitmq sender error")
	mockCommentDo.EXPECT().InsertComment(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 插入数据库成功，但删除缓存失败，且消息发布成功mock
	mockCommentDo.EXPECT().InsertComment(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 插入数据库成功，删除缓存成功mock
	mockCommentDo.EXPECT().InsertComment(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.AddCommentReq
		err  error
	}{
		{
			name: "add_comment_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Add comment with empty param"),
		},
		{
			name: "add_comment_with_empty_filed",
			req:  &pb.AddCommentReq{VideoId: 0, UserId: 1, Content: ""},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Add comment with empty video_id or user_id or content"),
		},
		{
			name: "add_comment_with_database_error",
			req:  &pb.AddCommentReq{VideoId: 1, UserId: 1, Content: "test"},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "insert comment failed: %v", insertCommentError),
		},
		{
			name: "add_comment_with_rabbit_error",
			req:  &pb.AddCommentReq{VideoId: 1, UserId: 1, Content: "test"},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "publish video comment count message failed: %v", senderError),
		},
		{
			name: "add_comment_with_redis_error",
			req:  &pb.AddCommentReq{VideoId: 1, UserId: 1, Content: "test"},
			err:  nil,
		},
		{
			name: "add_comment_success",
			req:  &pb.AddCommentReq{VideoId: 1, UserId: 1, Content: "test"},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := addCommentLogic.AddComment(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
