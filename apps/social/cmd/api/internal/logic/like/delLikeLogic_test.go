package like_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/social/cmd/api/internal/logic/like"
	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"
	socialMock "go-zero-douyin/apps/social/cmd/rpc/mock"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
	videoMock "go-zero-douyin/apps/video/cmd/rpc/mock"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"go-zero-douyin/mock"
	"testing"
)

func TestDelLikeLogic_DelLike(t *testing.T) {
	ctl := gomock.NewController(t)

	mockValidator := mock.NewMockValidator(ctl)

	mockVideoRpc := videoMock.NewMockVideo(ctl)

	mockSocialRpc := socialMock.NewMockSocial(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, VideoRpc: mockVideoRpc, SocialRpc: mockSocialRpc}

	delLikeLogic := like.NewDelLikeLogic(context.Background(), serviceContext)

	// 参数校验失败mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// 参数校验成功，但VideoRpc.GetVideoById失败mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	videoRpcError := errors.New("VideoRpc.GetVideoById error")
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(nil, videoRpcError)

	// 参数校验成功，且VideoRpc.GetVideoById调用成功，SocialRpc.VideoUnlike失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(&pb.GetVideoByIdResp{}, nil)
	socialRpcError := errors.New("SocialRpc.VideoUnlike error")
	mockSocialRpc.EXPECT().VideoUnlike(gomock.Any(), gomock.Any()).Return(nil, socialRpcError)

	// 参数校验成功，且VideoRpc.GetVideoById调用成功，SocialRpc.VideoUnlike成功的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(&pb.GetVideoByIdResp{}, nil)
	mockSocialRpc.EXPECT().VideoUnlike(gomock.Any(), gomock.Any()).Return(&socialPb.VideoUnlikeResp{}, nil)

	testCases := []struct {
		name string
		req  *types.VideoUnlikeReq
		err  error
	}{
		{
			name: "video_unlike_with_validate_error",
			req:  &types.VideoUnlikeReq{VideoId: 1},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "video_unlike_with_video_rpc_error",
			req:  &types.VideoUnlikeReq{VideoId: 1},
			err:  errors.Wrapf(videoRpcError, "req: %v", &types.VideoLikeReq{VideoId: 1}),
		},
		{
			name: "video_unlike_with_social_rpc_error",
			req:  &types.VideoUnlikeReq{VideoId: 1},
			err:  errors.Wrapf(socialRpcError, "req: %v", &types.VideoLikeReq{VideoId: 1}),
		},
		{
			name: "video_unlike_success",
			req:  &types.VideoUnlikeReq{VideoId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := delLikeLogic.DelLike(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
