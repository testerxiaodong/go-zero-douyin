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

func TestGetVideoLikedCountByVideoIdLogic_GetVideoLikedCountByVideoId(t *testing.T) {
	ctl := gomock.NewController(t)

	mockLikeDo := mock.NewMockLikeDo(ctl)

	mockRedis := globalMock.NewMockRedisCache(ctl)

	utils.IgnoreGo()
	defer utils.RecoverGo()

	serviceContext := &svc.ServiceContext{LikeDo: mockLikeDo, Redis: mockRedis, SingleFlight: syncx.NewSingleFlight()}

	getVideoLikedCountByVideoIdLogic := logic.NewGetVideoLikedCountByVideoIdLogic(context.Background(), serviceContext)

	// redis中有数据mock
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
	mockRedis.EXPECT().Scard(gomock.Any(), gomock.Any()).Return(int64(3), nil)
	mockRedis.EXPECT().Expire(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// redis中没有数据，查询数据库失败的mock
	dbError := errors.New("LikeDo GetVideoLikedCount error")
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
	mockLikeDo.EXPECT().GetVideoLikedCount(gomock.Any(), gomock.Any()).Return(int64(0), dbError)

	// redis中没有数据，查询数据库成功的mock
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
	mockLikeDo.EXPECT().GetVideoLikedCount(gomock.Any(), gomock.Any()).Return(int64(3), nil)

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
			name: "get_video_liked_count_by_video_id_with_redis",
			req:  &pb.GetVideoLikedCountByVideoIdReq{VideoId: 10},
			err:  nil,
		},
		{
			name: "get_video_liked_count_by_video_id_with_database_error",
			req:  &pb.GetVideoLikedCountByVideoIdReq{VideoId: 10},
			err:  dbError,
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
