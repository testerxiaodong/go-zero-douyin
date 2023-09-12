package video_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/video/cmd/api/internal/logic/video"
	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/mock"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"testing"
)

func TestGetAllTagLogic_GetAllTag(t *testing.T) {
	ctl := gomock.NewController(t)

	mockVideoRpc := mock.NewMockVideo(ctl)

	serviceContext := &svc.ServiceContext{VideoRpc: mockVideoRpc}

	getAllTagLogic := video.NewGetAllTagLogic(context.Background(), serviceContext)

	// videoRpc调用失败的mock
	videoRpcError := errors.New("videoRpc.GetAllTag error")
	mockVideoRpc.EXPECT().GetAllTag(gomock.Any(), gomock.Any()).Return(nil, videoRpcError)

	// videoRpc调用成功的mock
	expectedValue := &pb.TagInfo{Id: utils.NewRandomInt64(1, 10), Name: utils.NewRandomString(10)}
	mockVideoRpc.EXPECT().GetAllTag(gomock.Any(), gomock.Any()).
		Return(&pb.GetAllTagResp{Tags: []*pb.TagInfo{expectedValue}}, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "get_all_tag_with_video_rpc_error",
			err:  videoRpcError,
		},
		{
			name: "get_all_tag_success",
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := getAllTagLogic.GetAllTag()
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
				assert.Equal(t, len(resp.Tags), 1)
				assert.Equal(t, resp.Tags[0].Id, expectedValue.Id)
				assert.Equal(t, resp.Tags[0].Name, expectedValue.Name)
			}
		})
	}
}
