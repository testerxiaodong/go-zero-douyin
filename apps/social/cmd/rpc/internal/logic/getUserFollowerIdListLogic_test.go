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

func TestGetUserFollowerIdListLogic_GetUserFollowerIdList(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockFollowDo := mock.NewMockfollowModel(ctl)
	serviceContext := &svc.ServiceContext{FollowModel: mockFollowDo}
	getUserFollowerIdListLogic := logic.NewGetUserFollowerIdListLogic(context.Background(), serviceContext)

	// 查询数据库失败的mock
	dbError := errors.New("FollowDo GetUserFollowerIdList error")
	mockFollowDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockFollowDo.EXPECT().FindPageListByPage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, dbError)

	// 查询数据库成功，且没有数据的mock
	mockFollowDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockFollowDo.EXPECT().FindPageListByPage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*model.Follow{}, nil)

	// 查询数据库成功，且有数据的mock
	result := []*model.Follow{NewFollow(), NewFollow()}
	mockFollowDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockFollowDo.EXPECT().FindPageListByPage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(result, nil)

	// 表格驱动测试
	req := &pb.GetUserFollowerIdListReq{
		UserId:   utils.NewRandomInt64(1, 10),
		Page:     utils.NewRandomInt64(1, 5),
		PageSize: utils.NewRandomInt64(1, 10),
	}
	testCases := []struct {
		name string
		req  *pb.GetUserFollowerIdListReq
		err  error
	}{
		{
			name: "get_user_follower_id_list_with_empty_param",
			req:  nil,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follower id list with empty param"),
		},
		{
			name: "get_user_follower_id_list_with_empty_video_id",
			req:  &pb.GetUserFollowerIdListReq{UserId: 0},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follower id list with empty user_id"),
		},
		{
			name: "get_user_follower_id_list_over_page",
			req:  &pb.GetUserFollowerIdListReq{UserId: 10, Page: 6},
			err:  errors.Wrap(xerr.NewErrMsg("系统不允许超过五页"), "粉丝列表业务校验"),
		},
		{
			name: "get_user_follower_id_list_with_database_error",
			req:  req,
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"查询用户粉丝列表失败, err: %v, user_id: %d", dbError, req.GetUserId()),
		},
		{
			name: "get_user_follower_id_list_with_no_database_record",
			req:  req,
			err:  nil,
		},
		{
			name: "get_user_follower_id_list_with_database_record",
			req:  req,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getUserFollowerIdListLogic.GetUserFollowerIdList(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}

func NewFollow() *model.Follow {
	return &model.Follow{
		Id:         utils.NewRandomInt64(1, 10),
		UserId:     utils.NewRandomInt64(1, 10),
		FollowerId: utils.NewRandomInt64(1, 10),
		Status:     1,
	}
}
