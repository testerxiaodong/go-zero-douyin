package logic_test

import (
	"context"
	"github.com/Masterminds/squirrel"
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
	defer ctl.Finish()
	mockCommentDo := mock.NewMockcommentModel(ctl)
	serviceContext := &svc.ServiceContext{CommentModel: mockCommentDo}
	getVideoCommentListByIdLogic := logic.NewGetVideoCommentListByIdLogic(context.Background(), serviceContext)

	// 查询数据库失败mock
	dbError := errors.New("CommentDo.GetCommentListByVideoId error")
	mockCommentDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockCommentDo.EXPECT().FindPageListByPageWithTotal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, int64(0), dbError)

	// 查询数据库成功，且没有数据的mock
	mockCommentDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockCommentDo.EXPECT().FindPageListByPageWithTotal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*model.Comment{}, int64(0), nil)

	// 查询数据库成功，且有数据的mock
	expectedResp := []*model.Comment{&model.Comment{Id: 1, VideoId: 1, UserId: 1, Content: "test"},
		&model.Comment{Id: 2, VideoId: 2, UserId: 2, Content: "test"}}
	mockCommentDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockCommentDo.EXPECT().FindPageListByPageWithTotal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(expectedResp, int64(2), nil)

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
			name: "get_video_comment_list_by_id_with_no_record",
			req:  &pb.GetCommentListByIdReq{Id: 10},
			err:  nil,
		},
		{
			name: "get_video_comment_list_by_id_with_record",
			req:  &pb.GetCommentListByIdReq{Id: 10},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getVideoCommentListByIdLogic.GetVideoCommentListById(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
