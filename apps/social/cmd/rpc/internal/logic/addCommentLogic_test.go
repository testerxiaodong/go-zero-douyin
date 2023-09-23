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
	"testing"
)

func TestAddCommentLogic_AddComment(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockCommentDo := mock.NewMockcommentModel(ctl)
	mockCommentCountDo := mock.NewMockcommentCountModel(ctl)
	serviceContext := &svc.ServiceContext{CommentModel: mockCommentDo, CommentCountModel: mockCommentCountDo}
	addCommentLogic := logic.NewAddCommentLogic(context.Background(), serviceContext)

	// 事务失败的mock
	TransError := errors.New("CommentDo.Trans error")
	mockCommentDo.EXPECT().Trans(gomock.Any(), gomock.Any()).Return(TransError)

	// 事务成功的mock
	mockCommentDo.EXPECT().Trans(gomock.Any(), gomock.Any()).Return(nil)
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
			name: "add_comment_with_trans_error",
			req:  &pb.AddCommentReq{VideoId: 1, UserId: 1, Content: "test"},
			err:  TransError,
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
