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

func TestDelCommentLogic_DelComment(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockCommentDo := mock.NewMockcommentModel(ctl)
	serviceContext := &svc.ServiceContext{CommentModel: mockCommentDo}
	delCommentLogic := logic.NewDelCommentLogic(context.Background(), serviceContext)

	// 查询数据库失败的mock
	searchDatabaseError := errors.New("CommentDo.GetCommentById error")
	mockCommentDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, searchDatabaseError)

	// 查询数据库成功，但数据不存在的mock
	mockCommentDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, model.ErrNotFound)

	// 查询成功，数据存在，但非该用户发布的mock
	mockCommentDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(&model.Comment{UserId: 2}, nil)

	// 查询成功，数据存在，是该用户发布，但删除评论失败的mock
	transError := errors.New("CommentDo.Trans error")
	mockCommentDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(&model.Comment{UserId: 1}, nil)
	mockCommentDo.EXPECT().Trans(gomock.Any(), gomock.Any()).
		Return(transError)

	// 成功的mock
	mockCommentDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(&model.Comment{UserId: 1}, nil)
	mockCommentDo.EXPECT().Trans(gomock.Any(), gomock.Any()).
		Return(nil)

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
			err:  errors.Wrapf(xerr.NewErrMsg("评论不存在"), "comment not found, id: %d", 1),
		},
		{
			name: "del_comment_with_owner_error",
			req:  &pb.DelCommentReq{CommentId: 1, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("评论非该用户发布，无法删除"), "comment_id: %d", 1),
		},
		{
			name: "del_comment_with_trans_error",
			req:  &pb.DelCommentReq{CommentId: 1, UserId: 1},
			err:  transError,
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
