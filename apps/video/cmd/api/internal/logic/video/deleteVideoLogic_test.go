package video_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/video/cmd/api/internal/logic/video"
	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"
	"go-zero-douyin/apps/video/cmd/rpc/mock"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	gloablMock "go-zero-douyin/mock"
	"testing"
)

func TestDeleteVideoLogic_DeleteVideo(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	// 构造需要mock的接口
	mockVideoRpc := mock.NewMockVideo(ctl)
	mockValidator := gloablMock.NewMockValidator(ctl)

	// 创建deleteVideoLogic对象
	serviceContext := &svc.ServiceContext{Validator: mockValidator, VideoRpc: mockVideoRpc}
	videoDeleteLogic := video.NewDeleteVideoLogic(context.Background(), serviceContext)

	// mock具体的接口方法，实现测试逻辑

	// 参数校验失败mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// 参数校验成功，但videoRpc.DeleteVideo调用失败mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	rpcError := errors.New("videoRpc.DeleteVideo error")
	mockVideoRpc.EXPECT().DeleteVideo(gomock.Any(), gomock.Any()).Return(nil, rpcError)

	// 参数校验成功，且videoRpc.DeleteVideo调用成功的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().DeleteVideo(gomock.Any(), gomock.Any()).Return(&pb.DeleteVideoResp{}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *types.DeleteVideoReq
		err  error
	}{
		{
			name: "delete_video_with_validate_error",
			req:  &types.DeleteVideoReq{VideoId: 1},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "delete_video_with_video_rpc_error",
			req:  &types.DeleteVideoReq{VideoId: 1},
			err:  errors.Wrapf(rpcError, "req: %v", &types.DeleteVideoReq{VideoId: 1}),
		},
		{
			name: "delete_video_success",
			req:  &types.DeleteVideoReq{VideoId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := videoDeleteLogic.DeleteVideo(tc.req)
			if err != nil {
				assert.Equal(t, err.Error(), tc.err.Error())
			} else {
				assert.Equal(t, err, tc.err)
			}
		})
	}
}
