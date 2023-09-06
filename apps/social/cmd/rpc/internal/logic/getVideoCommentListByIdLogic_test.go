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

func TestGetVideoCommentListByIdLogic_GetVideoCommentListById(t *testing.T) {
	ctl := gomock.NewController(t)

	mockCommentDo := mock.NewMockCommentDo(ctl)

	serviceContext := &svc.ServiceContext{CommentDo: mockCommentDo}

	getVideoCommentListByIdLogic := logic.NewGetVideoCommentListByIdLogic(context.Background(), serviceContext)

	// 查询数据库失败mock
	dbError := errors.New("CommentDo.GetCommentListByVideoId error")
	mockCommentDo.EXPECT().GetCommentListByVideoId(gomock.Any(), gomock.Any()).Return(nil, dbError)

	// 查询数据库成功mock
	expectedResp := []*model.Comment{&model.Comment{ID: 1, VideoID: 1, UserID: 1, Content: "test"},
		&model.Comment{ID: 2, VideoID: 2, UserID: 2, Content: "test"}}
	mockCommentDo.EXPECT().GetCommentListByVideoId(gomock.Any(), gomock.Any()).
		Return(expectedResp, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.GetCommentListByIdReq
		err  error
	}{
		{
			name: "get_video_comment_list_by_id_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get video comment list by id with empty param"),
		},
		{
			name: "get_video_comment_list_by_id_with_empty_video_id",
			req:  &pb.GetCommentListByIdReq{Id: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get video comment list by id with empty id"),
		},
		{
			name: "get_video_comment_list_by_id_with_database_error",
			req:  &pb.GetCommentListByIdReq{Id: 10},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "db get comment list by video_id failed: %v", dbError),
		},
		{
			name: "get_video_comment_list_by_id_success",
			req:  &pb.GetCommentListByIdReq{Id: 10},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := getVideoCommentListByIdLogic.GetVideoCommentListById(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
				assert.Equal(t, len(resp.GetComments()), 2)
			}
		})
	}
}
