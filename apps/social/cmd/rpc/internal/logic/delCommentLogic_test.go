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

func TestDelCommentLogic_DelComment(t *testing.T) {
	ctl := gomock.NewController(t)

	mockCommentDo := mock.NewMockCommentDo(ctl)

	mockSender := gloabelMock.NewMockSender(ctl)

	mockRedis := gloabelMock.NewMockRedisCache(ctl)

	serviceContext := &svc.ServiceContext{CommentDo: mockCommentDo, Redis: mockRedis, Rabbit: mockSender}

	delCommentLogic := logic.NewDelCommentLogic(context.Background(), serviceContext)

	// 查询数据库失败的mock
	searchDatabaseError := errors.New("CommentDo.GetCommentById error")
	mockCommentDo.EXPECT().GetCommentById(gomock.Any(), gomock.Any()).Return(nil, searchDatabaseError)

	// 查询数据库成功，但数据不存在的mock
	mockCommentDo.EXPECT().GetCommentById(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	// 查询成功，数据存在，但非该用户发布的mock
	mockCommentDo.EXPECT().GetCommentById(gomock.Any(), gomock.Any()).Return(&model.Comment{UserID: 2}, nil)

	// 查询成功，数据存在，是该用户发布，但删除评论失败的mock
	deleteDatabaseError := errors.New("CommentDo.DeleteComment error")
	mockCommentDo.EXPECT().GetCommentById(gomock.Any(), gomock.Any()).Return(&model.Comment{UserID: 1}, nil)
	mockCommentDo.EXPECT().DeleteComment(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{RowsAffected: 0}, deleteDatabaseError)

	// 查询成功，数据存在，是该用户发布，删除评论成功，redis删除缓存失败，且rabbitmq发送消息失败的mock
	redisError := errors.New("redis delete error")
	senderError := errors.New("rabbitmq sender error")
	mockCommentDo.EXPECT().GetCommentById(gomock.Any(), gomock.Any()).Return(&model.Comment{UserID: 1}, nil)
	mockCommentDo.EXPECT().DeleteComment(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{RowsAffected: 1}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 查询成功，数据存在，是该用户发布，删除评论成功，redis删除缓存失败，且rabbitmq发送消息成功，但es消息发送失败的mock
	mockCommentDo.EXPECT().GetCommentById(gomock.Any(), gomock.Any()).Return(&model.Comment{UserID: 1}, nil)
	mockCommentDo.EXPECT().DeleteComment(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{RowsAffected: 1}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisError)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 查询成功，数据存在，是该用户发布，删除评论成功，redis删除缓存成功的mock
	mockCommentDo.EXPECT().GetCommentById(gomock.Any(), gomock.Any()).Return(&model.Comment{UserID: 1}, nil)
	mockCommentDo.EXPECT().DeleteComment(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{RowsAffected: 1}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 成功的mock
	mockCommentDo.EXPECT().GetCommentById(gomock.Any(), gomock.Any()).Return(&model.Comment{UserID: 1}, nil)
	mockCommentDo.EXPECT().DeleteComment(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{RowsAffected: 1}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.DelCommentReq
		err  error
	}{
		{
			name: "del_comment_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del comment with empty param"),
		},
		{
			name: "del_comment_with_empty_filed",
			req:  &pb.DelCommentReq{CommentId: 0, UserId: 1},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del comment with empty user_id or comment_id"),
		},
		{
			name: "del_comment_with_database_search_error",
			req:  &pb.DelCommentReq{CommentId: 1, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "find comment by id failed, err: %v", searchDatabaseError),
		},
		{
			name: "del_comment_with_database_no_record",
			req:  &pb.DelCommentReq{CommentId: 1, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "comment not found, id: %d", 1),
		},
		{
			name: "del_comment_with_owner_error",
			req:  &pb.DelCommentReq{CommentId: 1, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("评论非该用户发布，无法删除"), "comment_id: %d", 1),
		},
		{
			name: "del_comment_with_database_delete_error",
			req:  &pb.DelCommentReq{CommentId: 1, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "del comment failed, err: %v", deleteDatabaseError),
		},
		{
			name: "del_comment_with_rabbit_error",
			req:  &pb.DelCommentReq{CommentId: 1, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "publish video comment count message failed: %v", senderError),
		},
		{
			name: "del_comment_with_redis_error",
			req:  &pb.DelCommentReq{CommentId: 1, UserId: 1},
			err: errors.Wrapf(xerr.NewErrCode(xerr.RPC_UPDATE_ERR),
				"req: %v, err: %v", &pb.DelCommentReq{CommentId: 1, UserId: 1}, senderError),
		},
		{
			name: "del_comment_with_es_sender_error",
			req:  &pb.DelCommentReq{CommentId: 1, UserId: 1},
			err: errors.Wrapf(xerr.NewErrCode(xerr.RPC_UPDATE_ERR),
				"req: %v, err: %v", &pb.DelCommentReq{CommentId: 1, UserId: 1}, senderError),
		},
		{
			name: "del_comment_success",
			req:  &pb.DelCommentReq{CommentId: 1, UserId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := delCommentLogic.DelComment(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
