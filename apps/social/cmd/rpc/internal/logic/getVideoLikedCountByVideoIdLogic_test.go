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

func TestGetVideoLikedCountByVideoIdLogic_GetVideoLikedCountByVideoId(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockLikeDo := mock.NewMocklikeCountModel(ctl)
	serviceContext := &svc.ServiceContext{LikeCountModel: mockLikeDo}
	getVideoLikedCountByVideoIdLogic := logic.NewGetVideoLikedCountByVideoIdLogic(context.Background(), serviceContext)

	// 查询数据库失败的mock
	dbError := errors.New("LikeDo GetVideoLikedCount error")
	mockLikeDo.EXPECT().FindOneByVideoId(gomock.Any(), gomock.Any()).Return(nil, dbError)

	// 查询数据库成功，但没有数据的mock
	mockLikeDo.EXPECT().FindOneByVideoId(gomock.Any(), gomock.Any()).Return(nil, model.ErrNotFound)

	// 查询数据库成功，有数据的mock
	mockLikeDo.EXPECT().FindOneByVideoId(gomock.Any(), gomock.Any()).Return(&model.LikeCount{LikeCount: 1}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.GetVideoLikedCountByVideoIdReq
		err  error
	}{
		{
			name: "get_video_liked_count_by_video_id_with_empty_param",
			req:  nil,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video like count with empty param"),
		},
		{
			name: "get_video_liked_count_by_video_id_with_empty_video_id",
			req:  &pb.GetVideoLikedCountByVideoIdReq{VideoId: 0},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video like count with empty video_id"),
		},
		{
			name: "get_video_liked_count_by_video_id_with_database_error",
			req:  &pb.GetVideoLikedCountByVideoIdReq{VideoId: 10},
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"查询视频点赞数失败, err: %v, video_id: %d", dbError, 10),
		},
		{
			name: "get_video_liked_count_by_video_id_with_no_record",
			req:  &pb.GetVideoLikedCountByVideoIdReq{VideoId: 10},
			err:  nil,
		},
		{
			name: "get_video_liked_count_by_video_id_with_database_record",
			req:  &pb.GetVideoLikedCountByVideoIdReq{VideoId: 10},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getVideoLikedCountByVideoIdLogic.GetVideoLikedCountByVideoId(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
