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

func TestSearchLogic_Search(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockValidator := globalMock.NewMockValidator(ctl)
	mockVideoRpc := mock.NewMockVideo(ctl)
	serviceContext := &svc.ServiceContext{VideoRpc: mockVideoRpc, Validator: mockValidator}
	searchLogic := video.NewSearchLogic(context.Background(), serviceContext)

	// 参数校验失败的mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// VideoRpc.GetVideoByKeyword失败的mock
	videoRpcError := errors.New("VideoRpc.GetVideoByKeyword error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().GetVideoByKeyword(gomock.Any(), gomock.Any()).
		Return(nil, videoRpcError)

	// VideoRpc.GetVideoByKeyword成功但没数据的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().GetVideoByKeyword(gomock.Any(), gomock.Any()).
		Return(&pb.GetVideoByKeywordResp{}, nil)

	// VideoRpc.GetVideoByKeyword成功且有数据的mock
	expectedValue := &pb.VideoDetailInfo{
		Id:    utils.NewRandomInt64(1, 10),
		Title: utils.NewRandomString(10),
	}
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().GetVideoByKeyword(gomock.Any(), gomock.Any()).
		Return(&pb.GetVideoByKeywordResp{Videos: []*pb.VideoDetailInfo{expectedValue}}, nil)

	// 表格驱动测试
	req := &types.SearchVideoReq{Keyword: "test"}
	testCases := []struct {
		name string
		req  *types.SearchVideoReq
		err  error
	}{
		{
			name: "search_video_with_validate_error",
			req:  req,
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "search_video_with_video_rpc_error",
			req:  req,
			err:  errors.Wrapf(videoRpcError, "VideoRpc.GetVideoByKeyword失败, req: %v", req),
		},
		{
			name: "search_video_with_no_data",
			req:  req,
			err:  nil,
		},
		{
			name: "search_video_success",
			req:  req,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := searchLogic.Search(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
