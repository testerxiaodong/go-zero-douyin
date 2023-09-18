package search_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/search/cmd/api/internal/logic/search"
	"go-zero-douyin/apps/search/cmd/api/internal/svc"
	"go-zero-douyin/apps/search/cmd/api/internal/types"
	"go-zero-douyin/apps/search/cmd/rpc/mock"
	"go-zero-douyin/apps/search/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestCompleteVideoLogic_CompleteVideo(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockValidator := globalMock.NewMockValidator(ctl)
	mockSearchRpc := mock.NewMockSearch(ctl)
	serviceContext := &svc.ServiceContext{Validator: mockValidator, SearchRpc: mockSearchRpc}
	completeVideoLogic := search.NewCompleteVideoLogic(context.Background(), serviceContext)

	// 参数校验失败mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// searchRpc调用失败mock
	rpcError := errors.New("SearchRpc.CompleteVideo error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockSearchRpc.EXPECT().CompleteVideo(gomock.Any(), gomock.Any()).Return(nil, rpcError)

	// 成功的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockSearchRpc.EXPECT().CompleteVideo(gomock.Any(), gomock.Any()).
		Return(&pb.CompleteVideoResp{Suggestions: []string{"test"}}, nil)

	// 表格驱动测试
	req := &types.CompleteVideoReq{Input: "test"}
	testCases := []struct {
		name string
		req  *types.CompleteVideoReq
		err  error
	}{
		{
			name: "complete_video_with_validate_error",
			req:  req,
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "complete_video_with_rpc_error",
			req:  req,
			err:  errors.Wrapf(rpcError, "req: %v", req),
		},
		{
			name: "complete_video_success",
			req:  req,
			err:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := completeVideoLogic.CompleteVideo(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
