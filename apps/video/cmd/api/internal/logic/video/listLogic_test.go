package video_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	socialMock "go-zero-douyin/apps/social/cmd/rpc/mock"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
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

func TestListLogic_List(t *testing.T) {
	ctl := gomock.NewController(t)

	// 构造需要mock的接口
	mockVideoRpc := mock.NewMockVideo(ctl)
	mockSocialRpc := socialMock.NewMockSocial(ctl)
	mockValidator := gloablMock.NewMockValidator(ctl)

	// 创建deleteVideoLogic对象
	serviceContext := &svc.ServiceContext{Validator: mockValidator, VideoRpc: mockVideoRpc, SocialRpc: mockSocialRpc}
	listLogic := video.NewListLogic(context.Background(), serviceContext)

	// mock具体的接口方法，实现测试逻辑

	// 参数校验失败mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// 参数校验成功，但videoRpc.VideoFeed调用失败mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	videoRpcError := errors.New("videoRpc.UserVideoList error")
	mockVideoRpc.EXPECT().UserVideoList(gomock.Any(), gomock.Any()).Return(nil, videoRpcError)

	// 参数校验成功，且videoRpc.VideoFeed调用成功,但没有数据的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().UserVideoList(gomock.Any(), gomock.Any()).Return(&pb.UserVideoListResp{}, nil)

	// 参数校验成功，且videoRpc.VideoFeed调用成功,有两条数据的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	videoInfo := &pb.VideoInfo{Id: 1}
	socialRpcCommentCountError := errors.New("SocialRpc.GetCommentCountByVideoId error")
	socialRpcLikeCountError := errors.New("SocialRpc.GetVideoLikedCountByVideoId error")
	mockVideoRpc.EXPECT().UserVideoList(gomock.Any(), gomock.Any()).
		Return(&pb.UserVideoListResp{Videos: []*pb.VideoInfo{videoInfo, videoInfo}}, nil)

	mockSocialRpc.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).
		Return(nil, socialRpcCommentCountError)
	mockSocialRpc.EXPECT().GetVideoLikedCountByVideoId(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetVideoLikedCountByVideoIdResp{LikeCount: 1}, nil)

	mockSocialRpc.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetCommentCountByVideoIdResp{Count: 1}, nil)
	mockSocialRpc.EXPECT().GetVideoLikedCountByVideoId(gomock.Any(), gomock.Any()).
		Return(nil, socialRpcLikeCountError)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *types.UserVideoListReq
		err  error
	}{
		{
			name: "get_user_video_list_with_validate_error",
			req:  &types.UserVideoListReq{UserId: 1},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "get_user_video_list_with_video_rpc_error",
			req:  &types.UserVideoListReq{UserId: 1},
			err:  errors.Wrapf(videoRpcError, "req: %v", &types.UserVideoListReq{UserId: 1}),
		},
		{
			name: "get_user_video_list_with_no_data",
			req:  &types.UserVideoListReq{UserId: 1},
			err:  nil,
		},
		{
			name: "get_user_video_list_with_two_data",
			req:  &types.UserVideoListReq{UserId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := listLogic.List(tc.req)
			if err != nil {
				assert.Equal(t, err.Error(), tc.err.Error())
			} else {
				assert.Equal(t, err, tc.err)
			}
		})
	}

}
