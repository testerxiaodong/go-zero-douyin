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

func TestFollowUserLogic_FollowUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockFollowDo := mock.NewMockfollowModel(ctl)
	serviceContext := &svc.ServiceContext{FollowModel: mockFollowDo}
	followUserLogic := logic.NewFollowUserLogic(context.Background(), serviceContext)

	// 查询数据库失败mock
	dbSearchError := errors.New("FollowDo.GetFollowByFollowerIdAndUserId error")
	mockFollowDo.EXPECT().FindOneByUserIdFollowerIdIsDelete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, dbSearchError)

	// 查询成功，有记录，且关注状态为已关注的mock
	mockFollowDo.EXPECT().FindOneByUserIdFollowerIdIsDelete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.Follow{Status: xconst.FollowStateYes}, nil)

	// 查询成功，有记录，且关注状态为未关注，事务失败mock
	transError := errors.New("FollowDo.Trans error")
	mockFollowDo.EXPECT().FindOneByUserIdFollowerIdIsDelete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.Follow{Status: 0}, nil)
	mockFollowDo.EXPECT().Trans(gomock.Any(), gomock.Any()).Return(transError)

	// 查询成功，有记录，且关注状态为未关注，事务成功mock
	mockFollowDo.EXPECT().FindOneByUserIdFollowerIdIsDelete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&model.Follow{Status: 0}, nil)
	mockFollowDo.EXPECT().Trans(gomock.Any(), gomock.Any()).Return(nil)

	// 查询数据库成功，数据不存在，事务失败的mock
	mockFollowDo.EXPECT().FindOneByUserIdFollowerIdIsDelete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, model.ErrNotFound)
	mockFollowDo.EXPECT().Trans(gomock.Any(), gomock.Any()).Return(transError)

	// 查询数据库成功，数据不存在，事务成功的mock
	mockFollowDo.EXPECT().FindOneByUserIdFollowerIdIsDelete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, model.ErrNotFound)
	mockFollowDo.EXPECT().Trans(gomock.Any(), gomock.Any()).Return(nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.FollowUserReq
		err  error
	}{
		{
			name: "follow_user_with_empty_param",
			req:  nil,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "follow user with empty param"),
		},
		{
			name: "follow_user_with_empty_filed",
			req:  &pb.FollowUserReq{FollowerId: 0, UserId: 0},
			err: errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR),
				"follow user with empty follower_id or user_id"),
		},
		{
			name: "follow_user_self_error",
			req:  &pb.FollowUserReq{FollowerId: 1, UserId: 1},
			err: errors.Wrapf(xerr.NewErrMsg("不能对自己操作"),
				"req: %v", &pb.FollowUserReq{FollowerId: 1, UserId: 1}),
		},
		{
			name: "follow_user_with_database_search_error",
			req:  &pb.FollowUserReq{FollowerId: 2, UserId: 1},
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"search user is alreaddy follow user from db failed, err: %v, follower_id: %d user_id: %d",
				dbSearchError, 2, 1),
		},
		{
			name: "follow_user_with_already_follow",
			req:  &pb.FollowUserReq{FollowerId: 2, UserId: 1},
			err:  nil,
		},
		{
			name: "follow_user_with_have_record_trans_error",
			req:  &pb.FollowUserReq{FollowerId: 2, UserId: 1},
			err:  transError,
		},
		{
			name: "follow_user_have_record_success",
			req:  &pb.FollowUserReq{FollowerId: 2, UserId: 1},
			err:  nil,
		},
		{
			name: "follow_user_with_no_record_trans_error",
			req:  &pb.FollowUserReq{FollowerId: 2, UserId: 1},
			err:  transError,
		},
		{
			name: "follow_user_no_record_success",
			req:  &pb.FollowUserReq{FollowerId: 2, UserId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := followUserLogic.FollowUser(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
