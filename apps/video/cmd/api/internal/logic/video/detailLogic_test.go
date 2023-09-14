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
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestDetailLogic_Detail(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockValidator := globalMock.NewMockValidator(ctl)
	mockVideoRpc := mock.NewMockVideo(ctl)
	mockSocialRpc := socialMock.NewMockSocial(ctl)
	serviceContext := &svc.ServiceContext{Validator: mockValidator, VideoRpc: mockVideoRpc, SocialRpc: mockSocialRpc}
	detailLogic := video.NewDetailLogic(context.Background(), serviceContext)

	// 参数校验失败mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// videoRpc.GetVideoById失败的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	videoRpcError := errors.New("videoRpc.GetVideoById error")
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(nil, videoRpcError)

	// SocialRpc.GetCommentCountByVideoId失败的mock
	videoInfo := &pb.GetVideoByIdResp{
		Video: &pb.VideoInfo{
			Id:         utils.NewRandomInt64(1, 10),
			Title:      utils.NewRandomString(10),
			SectionId:  utils.NewRandomInt64(1, 10),
			Tags:       []string{"1", "2", "3"},
			OwnerId:    utils.NewRandomInt64(1, 10),
			PlayUrl:    utils.NewRandomString(10),
			CoverUrl:   utils.NewRandomString(10),
			CreateTime: utils.NewRandomInt64(1, 10),
			UpdateTime: utils.NewRandomInt64(1, 10),
		},
	}
	commentCountError := errors.New("SocialRpc.GetCommentCountByVideoId error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(videoInfo, nil)
	mockSocialRpc.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).Return(nil, commentCountError)

	// SocialRpc.GetVideoLikedCountByVideoId失败的mock
	likeCountError := errors.New("SocialRpc.GetVideoLikedCountByVideoId error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(videoInfo, nil)
	mockSocialRpc.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetCommentCountByVideoIdResp{Count: 1}, nil)
	mockSocialRpc.EXPECT().GetVideoLikedCountByVideoId(gomock.Any(), gomock.Any()).
		Return(nil, likeCountError)

	// 成功的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(videoInfo, nil)
	mockSocialRpc.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetCommentCountByVideoIdResp{Count: 1}, nil)
	mockSocialRpc.EXPECT().GetVideoLikedCountByVideoId(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetVideoLikedCountByVideoIdResp{LikeCount: 1}, nil)

	// 表格驱动测试
	req := &types.VideoDetailReq{VideoId: utils.NewRandomInt64(1, 10)}
	testCases := []struct {
		name string
		req  *types.VideoDetailReq
		err  error
	}{
		{
			name: "sync_video_to_es_with_validate_error",
			req:  req,
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "sync_video_to_es_with_video_rpc_error",
			req:  req,
			err:  errors.Wrapf(videoRpcError, "req: %v", req),
		},
		{
			name: "sync_video_to_es_with_comment_count_error",
			req:  req,
			err:  errors.Wrapf(commentCountError, "req: %v", req),
		},
		{
			name: "sync_video_to_es_with_like_count_error",
			req:  req,
			err:  errors.Wrapf(likeCountError, "req: %v", req),
		},
		{
			name: "sync_video_to_es_success",
			req:  req,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := detailLogic.Detail(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
