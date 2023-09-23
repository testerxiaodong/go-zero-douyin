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
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"testing"
)

func TestGetUserLikeVideoIdListLogic_GetUserLikeVideoIdList(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockLikeDo := mock.NewMocklikeModel(ctl)
	serviceContext := &svc.ServiceContext{LikeModel: mockLikeDo}
	getUserLikeVideoIdListLogic := logic.NewGetUserLikeVideoIdListLogic(context.Background(), serviceContext)

	// 查询数据库失败的mock
	dbError := errors.New("FollowDo GetUserFollowerIdList error")
	mockLikeDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockLikeDo.EXPECT().FindPageListByPage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, dbError)

	// 查询数据库成功，但没有数据的mock
	mockLikeDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockLikeDo.EXPECT().FindPageListByPage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*model.Like{}, nil)

	// 查询数据库成功，有数据的mock
	result := []*model.Like{NewLike(), NewLike()}
	mockLikeDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockLikeDo.EXPECT().FindPageListByPage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(result, nil)

	// 表格驱动测试
	req := &pb.GetUserLikeVideoIdListReq{
		UserId:   utils.NewRandomInt64(1, 10),
		Page:     utils.NewRandomInt64(1, 10),
		PageSize: utils.NewRandomInt64(1, 10),
	}
	testCases := []struct {
		name string
		req  *pb.GetUserLikeVideoIdListReq
		err  error
	}{
		{
			name: "get_user_like_video_id_list_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user like video id list with empty param"),
		},
		{
			name: "get_user_like_video_id_list_with_empty_video_id",
			req:  &pb.GetUserLikeVideoIdListReq{UserId: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user like video id list with empty user_id"),
		},
		{
			name: "get_user_like_video_id_list_with_database_error",
			req:  req,
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"查询用户点赞视频列表失败, err: %v, user_id: %d", dbError, req.GetUserId()),
		},
		{
			name: "get_user_like_video_id_list_with_no_database_record",
			req:  req,
			err:  nil,
		},
		{
			name: "get_user_like_video_id_list_with_database_record",
			req:  req,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getUserLikeVideoIdListLogic.GetUserLikeVideoIdList(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}

func NewLike() *model.Like {
	return &model.Like{
		Id:      utils.NewRandomInt64(1, 10),
		UserId:  utils.NewRandomInt64(1, 10),
		Version: utils.NewRandomInt64(1, 10),
		Status:  utils.NewRandomInt64(0, 1),
	}
}
