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

func TestGetCommentCountByVideoIdLogic_GetCommentCountByVideoId(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockCommentDo := mock.NewMockcommentCountModel(ctl)
	serviceContext := &svc.ServiceContext{CommentCountModel: mockCommentDo}
	getCommentCountByVideoIdLogic := logic.NewGetCommentCountByVideoIdLogic(context.Background(), serviceContext)

	// 查询数据库失败的mock
	dbError := errors.New("CommentDo GetCommentCountByVideoId error")
	mockCommentDo.EXPECT().FindOneByVideoId(gomock.Any(), gomock.Any()).Return(nil, dbError)

	// 查询数据库成功，但是记录不存在的mock
	mockCommentDo.EXPECT().FindOneByVideoId(gomock.Any(), gomock.Any()).Return(nil, model.ErrNotFound)

	// 查询数据库成功，且记录存在的mock
	mockCommentDo.EXPECT().FindOneByVideoId(gomock.Any(), gomock.Any()).
		Return(&model.CommentCount{CommentCount: 1}, nil)
	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.GetCommentCountByVideoIdReq
		err  error
	}{
		{
			name: "get_comment_count_by_video_id_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video comment count with empty param"),
		},
		{
			name: "get_comment_count_by_video_id_with_empty_video_id",
			req:  &pb.GetCommentCountByVideoIdReq{VideoId: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video comment count with empty video_id"),
		},
		{
			name: "get_comment_count_by_video_id_with_database_error",
			req:  &pb.GetCommentCountByVideoIdReq{VideoId: 1},
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"查询视频评论数失败, err: %v, video_id: %d", dbError, 1),
		},
		{
			name: "get_comment_count_by_video_id_with_no_database_record",
			req:  &pb.GetCommentCountByVideoIdReq{VideoId: 1},
			err:  nil,
		},
		{
			name: "get_comment_count_by_video_id_success",
			req:  &pb.GetCommentCountByVideoIdReq{VideoId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getCommentCountByVideoIdLogic.GetCommentCountByVideoId(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
