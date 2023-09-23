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

func TestUnfollowUserLogic_UnfollowUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockFollowDo := mock.NewMockfollowModel(ctl)
	serviceContext := &svc.ServiceContext{FollowModel: mockFollowDo}
	unfollowUserLogic := logic.NewUnfollowUserLogic(context.Background(), serviceContext)

	// 查询数据库失败mock
	dbSearchError := errors.New("FollowDo.GetFollowByFollowerIdAndUserId error")
	mockFollowDo.EXPECT().FindOneByUserIdFollowerId(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, dbSearchError)

	// 查询数据库成功，数据不存在的mock
	mockFollowDo.EXPECT().FindOneByUserIdFollowerId(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, model.ErrNotFound)

	// 查询数据库成功，有记录，且为未关注
	mockFollowDo.EXPECT().FindOneByUserIdFollowerId(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.Follow{Status: xconst.FollowStateNo}, nil)

	// 查询数据库成功，有记录，且为已关注，事务失败的mock
	transError := errors.New("FollowModel.Trans error")
	mockFollowDo.EXPECT().FindOneByUserIdFollowerId(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.Follow{Status: xconst.FollowStateYes}, nil)
	mockFollowDo.EXPECT().Trans(gomock.Any(), gomock.Any()).Return(transError)

	// 查询数据库成功，有记录，且为已关注，事务成功的mock
	mockFollowDo.EXPECT().FindOneByUserIdFollowerId(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.Follow{Status: xconst.FollowStateYes}, nil)
	mockFollowDo.EXPECT().Trans(gomock.Any(), gomock.Any()).Return(nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.UnfollowUserReq
		err  error
	}{
		{
			name: "unfollow_user_with_empty_param",
			req:  nil,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unfollow user with empty param"),
		},
		{
			name: "unfollow_user_with_empty_filed",
			req:  &pb.UnfollowUserReq{FollowerId: 0, UserId: 0},
			err: errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR),
				"unfollow user with empty follower_id or user_id"),
		},
		{
			name: "unfollow_user_self_error",
			req:  &pb.UnfollowUserReq{FollowerId: 1, UserId: 1},
			err: errors.Wrapf(xerr.NewErrMsg("不能对自己操作"),
				"req: %v", &pb.FollowUserReq{FollowerId: 1, UserId: 1}),
		},
		{
			name: "unfollow_user_with_database_search_error",
			req:  &pb.UnfollowUserReq{FollowerId: 2, UserId: 1},
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"查询关注记录失败, err: %v, user_id: %v, follower_id: %d", dbSearchError, 1, 2),
		},
		{
			name: "unfollow_user_with_database_no_record",
			req:  &pb.UnfollowUserReq{FollowerId: 2, UserId: 1},
			err: errors.Wrapf(xerr.NewErrMsg("未关注过该用户"),
				"user_id: %v, follower_id: %d", 1, 2),
		},
		{
			name: "unfollow_user_with_unfollow_record",
			req:  &pb.UnfollowUserReq{FollowerId: 2, UserId: 1},
			err:  nil,
		},
		{
			name: "unfollow_user_with_follow_record_trans_error",
			req:  &pb.UnfollowUserReq{FollowerId: 2, UserId: 1},
			err:  transError,
		},
		{
			name: "unfollow_user_with_follow_record_success",
			req:  &pb.UnfollowUserReq{FollowerId: 2, UserId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := unfollowUserLogic.UnfollowUser(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
