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
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestDelSectionLogic_DelSection(t *testing.T) {
	ctl := gomock.NewController(t)

	mockValidator := globalMock.NewMockValidator(ctl)

	mockVideoRpc := mock.NewMockVideo(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, VideoRpc: mockVideoRpc}

	delSectionLogic := video.NewDelSectionLogic(context.Background(), serviceContext)

	// 参数校验失败的mock
	validatorError := "validate error"
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validatorError)

	// 参数校验成功，但videoRpc调用失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	videoRpcError := errors.New("VideoRpc.DelSection error")
	mockVideoRpc.EXPECT().DelSection(gomock.Any(), gomock.Any()).Return(nil, videoRpcError)

	// 参数校验成功，且videoRpc调用成功的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().DelSection(gomock.Any(), gomock.Any()).Return(&pb.DelSectionResp{}, nil)

	param := &types.DelSectionReq{Id: utils.NewRandomInt64(1, 10)}

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *types.DelSectionReq
		err  error
	}{
		{
			name: "del_section_with_validate_error",
			req:  param,
			err:  xerr.NewErrMsg(validatorError),
		},
		{
			name: "del_section_with_video_rpc_error",
			req:  param,
			err:  errors.Wrapf(videoRpcError, "req: %v", param),
		},
		{
			name: "del_section_success",
			req:  param,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := delSectionLogic.DelSection(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
