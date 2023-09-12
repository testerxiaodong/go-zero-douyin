package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/syncx"
	"go-zero-douyin/apps/social/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/mock"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestGetCommentCountByVideoIdLogic_GetCommentCountByVideoId(t *testing.T) {
	ctl := gomock.NewController(t)

	mockCommentDo := mock.NewMockCommentDo(ctl)

	mockRedis := globalMock.NewMockRedisCache(ctl)

	utils.IgnoreGo()
	defer utils.RecoverGo()

	serviceContext := &svc.ServiceContext{CommentDo: mockCommentDo, Redis: mockRedis, SingleFlight: syncx.NewSingleFlight()}

	getCommentCountByVideoIdLogic := logic.NewGetCommentCountByVideoIdLogic(context.Background(), serviceContext)

	// redis中有数据mock
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
	mockRedis.EXPECT().Get(gomock.Any(), gomock.Any()).Return("1", nil)
	mockRedis.EXPECT().Expire(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// redis中没有数据，查询数据库失败的mock
	dbError := errors.New("CommentDo GetCommentCountByVideoId error")
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
	mockCommentDo.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).Return(int64(0), dbError)

	// redis中没有数据，查询数据库成功的mock
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
	mockCommentDo.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).Return(int64(10), nil)

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
			name: "get_comment_count_by_video_id_with_redis",
			req:  &pb.GetCommentCountByVideoIdReq{VideoId: 1},
			err:  nil,
		},
		{
			name: "get_comment_count_by_video_id_with_database_error",
			req:  &pb.GetCommentCountByVideoIdReq{VideoId: 1},
			err:  dbError,
		},
		{
			name: "get_comment_count_by_video_id_with_database_record",
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
