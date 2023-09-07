package comment_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/social/cmd/api/internal/logic/comment"
	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"
	socialMock "go-zero-douyin/apps/social/cmd/rpc/mock"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/apps/video/cmd/rpc/mock"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	gloablMock "go-zero-douyin/mock"
	"testing"
)

func TestAddCommentLogic_AddComment(t *testing.T) {
	ctl := gomock.NewController(t)

	mockValidator := gloablMock.NewMockValidator(ctl)

	mockVideoRpc := mock.NewMockVideo(ctl)

	mockSocialRpc := socialMock.NewMockSocial(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, VideoRpc: mockVideoRpc, SocialRpc: mockSocialRpc}

	addCommentLogic := comment.NewAddCommentLogic(context.Background(), serviceContext)

	// 参数校验失败mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// 参数校验成功，但VideoRpc.GetVideoById失败mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	videoRpcError := errors.New("VideoRpc.GetVideoById error")
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(nil, videoRpcError)

	// 参数校验成功，且VideoRpc.GetVideoById调用成功，SocialRpc.AddComment失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(&pb.GetVideoByIdResp{}, nil)
	socialRpcError := errors.New("SocialRpc.AddComment error")
	mockSocialRpc.EXPECT().AddComment(gomock.Any(), gomock.Any()).Return(nil, socialRpcError)

	// 参数校验成功，且VideoRpc.GetVideoById调用成功，SocialRpc.AddComment成功的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(&pb.GetVideoByIdResp{}, nil)
	mockSocialRpc.EXPECT().AddComment(gomock.Any(), gomock.Any()).Return(&socialPb.AddCommentResp{}, nil)

	testCases := []struct {
		name string
		req  *types.AddCommentReq
		err  error
	}{
		{
			name: "add_comment_with_validate_error",
			req:  &types.AddCommentReq{VideoId: 1, Content: "test"},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "add_comment_with_video_rpc_error",
			req:  &types.AddCommentReq{VideoId: 1, Content: "test"},
			err:  videoRpcError,
		},
		{
			name: "add_comment_with_social_rpc_error",
			req:  &types.AddCommentReq{VideoId: 1, Content: "test"},
			err:  socialRpcError,
		},
		{
			name: "add_comment_success",
			req:  &types.AddCommentReq{VideoId: 1, Content: "test"},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := addCommentLogic.AddComment(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
