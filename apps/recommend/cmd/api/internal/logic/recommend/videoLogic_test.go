package recommend_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/recommend/cmd/api/internal/logic/recommend"
	"go-zero-douyin/apps/recommend/cmd/api/internal/svc"
	"go-zero-douyin/apps/recommend/cmd/api/internal/types"
	"go-zero-douyin/apps/recommend/cmd/rpc/mock"
	"go-zero-douyin/apps/recommend/cmd/rpc/pb"
	socialMock "go-zero-douyin/apps/social/cmd/rpc/mock"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
	videoMock "go-zero-douyin/apps/video/cmd/rpc/mock"
	videoPb "go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestVideoLogic_Video(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockValidator := globalMock.NewMockValidator(ctl)
	mockRecommendRpc := mock.NewMockRecommend(ctl)
	mockVideoRpc := videoMock.NewMockVideo(ctl)
	mockSocialRpc := socialMock.NewMockSocial(ctl)
	serviceContext := &svc.ServiceContext{Validator: mockValidator, RecommendRpc: mockRecommendRpc,
		VideoRpc: mockVideoRpc, SocialRpc: mockSocialRpc}
	videoLogic := recommend.NewVideoLogic(context.Background(), serviceContext)

	// 参数校验失败的mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// RecommendRpc调用失败的mock
	recommendRpcError := errors.New("recommend rpc error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockRecommendRpc.EXPECT().VideoRecommendSection(gomock.Any(), gomock.Any()).Return(nil, recommendRpcError)

	// VideoRpc调用失败的mock
	videoRpcError := errors.New("video rpc error")
	recommendResult := []int64{utils.NewRandomInt64(1, 10)}
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockRecommendRpc.EXPECT().VideoRecommendSection(gomock.Any(), gomock.Any()).
		Return(&pb.VideoRecommendSectionResp{VideoIds: recommendResult}, nil)
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(nil, videoRpcError)
	mockSocialRpc.EXPECT().GetVideoLikedCountByVideoId(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetVideoLikedCountByVideoIdResp{LikeCount: 1}, nil)
	mockSocialRpc.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetCommentCountByVideoIdResp{Count: 1}, nil)

	// SocialRpc.GetVideoLikedCountByVideoId失败的mock
	video := &videoPb.VideoInfo{
		Id:        utils.NewRandomInt64(1, 10),
		Title:     utils.NewRandomString(10),
		TagIds:    utils.NewRandomString(10),
		SectionId: utils.NewRandomInt64(1, 10),
		OwnerId:   utils.NewRandomInt64(1, 10),
		OwnerName: utils.NewRandomString(10),
		PlayUrl:   utils.NewRandomString(10),
		CoverUrl:  utils.NewRandomString(10),
	}
	socialRpcLikeCountError := errors.New("socialRpc.GetVideoLikedCountByVideoId error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockRecommendRpc.EXPECT().VideoRecommendSection(gomock.Any(), gomock.Any()).
		Return(&pb.VideoRecommendSectionResp{VideoIds: recommendResult}, nil)
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).
		Return(&videoPb.GetVideoByIdResp{Video: video}, nil)
	mockSocialRpc.EXPECT().GetVideoLikedCountByVideoId(gomock.Any(), gomock.Any()).
		Return(nil, socialRpcLikeCountError)
	mockSocialRpc.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetCommentCountByVideoIdResp{Count: 1}, nil)

	// SocialRpc.GetCommentCountByVideoId失败的mock
	socialRpcCommentCountError := errors.New("socialRpc.GetCommentCountByVideoId error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockRecommendRpc.EXPECT().VideoRecommendSection(gomock.Any(), gomock.Any()).
		Return(&pb.VideoRecommendSectionResp{VideoIds: recommendResult}, nil)
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).
		Return(&videoPb.GetVideoByIdResp{Video: video}, nil)
	mockSocialRpc.EXPECT().GetVideoLikedCountByVideoId(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetVideoLikedCountByVideoIdResp{LikeCount: 1}, nil)
	mockSocialRpc.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).
		Return(nil, socialRpcCommentCountError)

	// 成功的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockRecommendRpc.EXPECT().VideoRecommendSection(gomock.Any(), gomock.Any()).
		Return(&pb.VideoRecommendSectionResp{VideoIds: recommendResult}, nil)
	mockVideoRpc.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).
		Return(&videoPb.GetVideoByIdResp{Video: video}, nil)
	mockSocialRpc.EXPECT().GetVideoLikedCountByVideoId(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetVideoLikedCountByVideoIdResp{LikeCount: 1}, nil)
	mockSocialRpc.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetCommentCountByVideoIdResp{Count: 1}, nil)

	// 表格驱动测试
	req := &types.VideoRecommendReq{
		SectionId: utils.NewRandomInt64(1, 10),
		Count:     utils.NewRandomInt64(1, 10),
	}
	testCases := []struct {
		name string
		req  *types.VideoRecommendReq
		err  error
	}{
		{
			name: "video_recommend_section_with_validate_error",
			req:  req,
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "video_recommend_section_with_recommend_rpc_error",
			req:  req,
			err:  errors.Wrapf(recommendRpcError, "req: %v", req),
		},
		{
			name: "video_recommend_section_with_video_rpc_error",
			req:  req,
			err:  errors.Wrapf(videoRpcError, "req: %v", req),
		},
		{
			name: "video_recommend_section_with_social_rpc_like_count_error",
			req:  req,
			err:  errors.Wrapf(socialRpcLikeCountError, "req: %v", req),
		},
		{
			name: "video_recommend_section_with_social_rpc_comment_count_error",
			req:  req,
			err:  errors.Wrapf(socialRpcCommentCountError, "req: %v", req),
		},
		{
			name: "video_recommend_section_with_success",
			req:  req,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := videoLogic.Video(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
