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
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"
	"testing"
)

func TestVideoUnlikeLogic_VideoUnlike(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockLikeDo := mock.NewMocklikeModel(ctl)
	serviceContext := &svc.ServiceContext{LikeModel: mockLikeDo}
	videoUnlikeLogic := logic.NewVideoUnlikeLogic(context.Background(), serviceContext)

	// 查询数据库失败mock
	dbSearchError := errors.New("LikeDo.GetLikeByVideoIdAndUserId error")
	mockLikeDo.EXPECT().FindOneByVideoIdUserId(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, dbSearchError)

	// 查询数据库成功，数据不存在的mock
	mockLikeDo.EXPECT().FindOneByVideoIdUserId(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, model.ErrNotFound)

	// 查询数据库成功，数据存在，且状态为已取消点赞的mock
	mockLikeDo.EXPECT().FindOneByVideoIdUserId(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.Like{Status: xconst.LikeStateNo}, nil)

	// 查询数据库成功，数据存在，且状态为已点赞，事务失败的mock
	transError := errors.New("LikeModel.Trans error")
	mockLikeDo.EXPECT().FindOneByVideoIdUserId(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.Like{Status: xconst.LikeStateYes}, nil)
	mockLikeDo.EXPECT().Trans(gomock.Any(), gomock.Any()).Return(transError)

	// 查询数据库成功，数据存在，且状态为已点赞，事务成功的mock
	mockLikeDo.EXPECT().FindOneByVideoIdUserId(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.Like{Status: xconst.LikeStateYes}, nil)
	mockLikeDo.EXPECT().Trans(gomock.Any(), gomock.Any()).Return(nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.VideoUnlikeReq
		err  error
	}{
		{
			name: "video_unlike_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unlike video with empty param"),
		},
		{
			name: "video_unlike_with_empty_filed",
			req:  &pb.VideoUnlikeReq{VideoId: 0, UserId: 0},
			err: errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR),
				"unlike video with empty video_id or user_id"),
		},
		{
			name: "video_unlike_with_database_search_error",
			req:  &pb.VideoUnlikeReq{VideoId: 2, UserId: 1},
			err: errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR),
				"find video is liked by user failed, err: %v", dbSearchError),
		},
		{
			name: "video_unlike_with_database_no_record",
			req:  &pb.VideoUnlikeReq{VideoId: 2, UserId: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("用户尚未点赞"), "video_id: %d", 2),
		},
		{
			name: "video_unlike_with_database_record_unlike",
			req:  &pb.VideoUnlikeReq{VideoId: 2, UserId: 1},
			err:  nil,
		},
		{
			name: "video_unlike_with_database_record_like_trans_error",
			req:  &pb.VideoUnlikeReq{VideoId: 2, UserId: 1},
			err:  transError,
		},
		{
			name: "video_unlike_with_database_record_like_trans_success",
			req:  &pb.VideoUnlikeReq{VideoId: 2, UserId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := videoUnlikeLogic.VideoUnlike(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
