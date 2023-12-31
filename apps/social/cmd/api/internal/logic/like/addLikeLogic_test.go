package like_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/social/cmd/api/internal/logic/like"
	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"
	videoMock "go-zero-douyin/apps/video/cmd/rpc/mock"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"go-zero-douyin/mock"
	"testing"
)

func TestAddLikeLogic_AddLike(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockValidator := mock.NewMockValidator(ctl)
	mockPusher := mock.NewMockPusher(ctl)
	mockVideoRpc := videoMock.NewMockVideo(ctl)
	serviceContext := &svc.ServiceContext{Validator: mockValidator, VideoRpc: mockVideoRpc, KqueueVideoLikeClient: mockPusher}
	addLikeLogic := like.NewAddLikeLogic(context.Background(), serviceContext)

	// 参数校验失败mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// 参数校验成功，但VideoRpc.GetVideoById失败mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	videoRpcError := errors.New("VideoRpc.GetVideoById error")
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(nil, videoRpcError)

	// 参数校验成功，且VideoRpc.GetVideoById调用成功，Pusher失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(&pb.GetVideoByIdResp{}, nil)
	pusherError := errors.New("SocialRpc.VideoLike error")
	mockPusher.EXPECT().Push(gomock.Any()).Return(pusherError)

	// 参数校验成功，且VideoRpc.GetVideoById调用成功，Pusher成功的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(&pb.GetVideoByIdResp{}, nil)
	mockPusher.EXPECT().Push(gomock.Any()).Return(nil)

	testCases := []struct {
		name string
		req  *types.VideoLikeReq
		err  error
	}{
		{
			name: "video_like_with_validate_error",
			req:  &types.VideoLikeReq{VideoId: 1},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "video_like_with_video_rpc_error",
			req:  &types.VideoLikeReq{VideoId: 1},
			err:  errors.Wrapf(videoRpcError, "req: %v", &types.VideoLikeReq{VideoId: 1}),
		},
		{
			name: "video_like_with_social_rpc_error",
			req:  &types.VideoLikeReq{VideoId: 1},
			err:  errors.Wrapf(pusherError, "req: %v", &types.VideoLikeReq{VideoId: 1}),
		},
		{
			name: "video_like_success",
			req:  &types.VideoLikeReq{VideoId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := addLikeLogic.AddLike(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
