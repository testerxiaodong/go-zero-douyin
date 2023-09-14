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
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestDeleteUserLogic_DeleteUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockElastic := globalMock.NewMockElasticService(ctl)
	serviceContext := &svc.ServiceContext{ElasticSearch: mockElastic}
	deleteUserLogic := logic.NewDeleteUserLogic(context.Background(), serviceContext)

	// ElasticSearch.DeleteDocument失败的mock
	esError := errors.New("es error")
	mockElastic.EXPECT().DeleteDocument(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, esError)

	// ElasticSearch.DeleteDocument成功的mock
	mockElastic.EXPECT().DeleteDocument(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&elastic.DeleteResponse{}, nil)

	// 表格驱动测试
	req := &pb.DeleteUserDocumentReq{Id: 1}
	testCases := []struct {
		name string
		req  *pb.DeleteUserDocumentReq
		err  error
	}{
		{
			name: "delete_user_with_nil_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "参数不能为nil"),
		},
		{
			name: "delete_user_with_empty_id",
			req:  &pb.DeleteUserDocumentReq{Id: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "用户id不能为空"),
		},
		{
			name: "delete_user_with_es_error",
			req:  req,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.RPC_DELETE_ERR), "删除用户文档失败, err: %v", esError),
		},
		{
			name: "delete_user_success",
			req:  req,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := deleteUserLogic.DeleteUser(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
