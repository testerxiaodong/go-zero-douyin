package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/search/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/search/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/search/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestSyncUserInfoLogic_SyncUserInfo(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockElastic := globalMock.NewMockElasticService(ctl)
	serviceContext := &svc.ServiceContext{ElasticSearch: mockElastic}
	syncUserInfoLogic := logic.NewSyncUserInfoLogic(context.Background(), serviceContext)

	// ElasticSearch.DeleteDocument失败的mock
	esError := errors.New("es error")
	mockElastic.EXPECT().CreateDocument(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, esError)

	// ElasticSearch.DeleteDocument成功的mock
	mockElastic.EXPECT().CreateDocument(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&elastic.IndexResponse{}, nil)

	// 表格驱动测试
	req := &pb.SyncUserInfoReq{User: &pb.User{
		Id:            utils.NewRandomInt64(1, 10),
		Username:      utils.NewRandomString(10),
		FollowerCount: utils.NewRandomInt64(1, 10),
		FollowCount:   utils.NewRandomInt64(1, 10),
	}}
	testCases := []struct {
		name string
		req  *pb.SyncUserInfoReq
		err  error
	}{
		{
			name: "sync_user_info_with_nil_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "参数不能为nil"),
		},
		{
			name: "sync_user_info_with_empty_id",
			req:  &pb.SyncUserInfoReq{User: &pb.User{Id: 0}},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "用户id不允许为空"),
		},
		{
			name: "sync_user_info_with_es_error",
			req:  req,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.RPC_UPDATE_ERR), "更新es用户文档失败, err: %v", esError),
		},
		{
			name: "sync_user_info_success",
			req:  req,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := syncUserInfoLogic.SyncUserInfo(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
