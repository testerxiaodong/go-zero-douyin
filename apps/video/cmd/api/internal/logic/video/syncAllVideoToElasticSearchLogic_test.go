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
	"go-zero-douyin/apps/video/cmd/rpc/mock"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"testing"
)

func TestSyncAllVideoToElasticsearchLogic_SyncAllVideoToElasticsearch(t *testing.T) {
	ctl := gomock.NewController(t)

	// 构造需要mock的接口
	mockVideoRpc := mock.NewMockVideo(ctl)
	mockSocialRpc := socialMock.NewMockSocial(ctl)

	// 创建deleteVideoLogic对象
	serviceContext := &svc.ServiceContext{VideoRpc: mockVideoRpc, SocialRpc: mockSocialRpc}
	syncAllVideoToElasticsearchLogic := video.NewSyncAllVideoToElasticsearchLogic(context.Background(), serviceContext)

	// mock具体的接口方法，实现测试逻辑

	// 参数校验成功，但videoRpc.VideoFeed调用失败mock
	videoRpcError := errors.New("videoRpc.UserVideoList error")
	mockVideoRpc.EXPECT().GetAllVideo(gomock.Any(), gomock.Any()).Return(nil, videoRpcError)

	// 参数校验成功，且videoRpc.VideoFeed调用成功,但没有数据的mock
	mockVideoRpc.EXPECT().GetAllVideo(gomock.Any(), gomock.Any()).Return(&pb.GetAllVideoResp{}, nil)

	// 参数校验成功，且videoRpc.VideoFeed调用成功,有两条数据的mock
	videoInfo := &pb.VideoInfo{Id: 1}
	socialRpcCommentCountError := errors.New("SocialRpc.GetCommentCountByVideoId error")
	socialRpcLikeCountError := errors.New("SocialRpc.GetVideoLikedCountByVideoId error")
	mockVideoRpc.EXPECT().GetAllVideo(gomock.Any(), gomock.Any()).
		Return(&pb.GetAllVideoResp{Videos: []*pb.VideoInfo{videoInfo, videoInfo}}, nil)

	videoRpcSyncVideoInfoToElasticsearchError := errors.New("SyncVideoInfoToElasticsearchError")

	mockSocialRpc.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).
		Return(nil, socialRpcCommentCountError)
	mockSocialRpc.EXPECT().GetVideoLikedCountByVideoId(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetVideoLikedCountByVideoIdResp{LikeCount: 1}, nil)
	mockVideoRpc.EXPECT().SyncVideoInfoToElasticsearch(gomock.Any(), gomock.Any()).
		Return(&pb.SyncVideoInfoToElasticsearchResp{}, nil)

	mockSocialRpc.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).
		Return(&socialPb.GetCommentCountByVideoIdResp{Count: 1}, nil)
	mockSocialRpc.EXPECT().GetVideoLikedCountByVideoId(gomock.Any(), gomock.Any()).
		Return(nil, socialRpcLikeCountError)
	mockVideoRpc.EXPECT().SyncVideoInfoToElasticsearch(gomock.Any(), gomock.Any()).
		Return(nil, videoRpcSyncVideoInfoToElasticsearchError)

	// 表格驱动测试
	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "sync_video_info_to_es_with_video_rpc_error",
			err:  errors.Wrap(videoRpcError, "获取所有视频信息失败"),
		},
		{
			name: "sync_video_info_to_es_with_no_data",
			err:  nil,
		},
		{
			name: "sync_video_info_to_es_with_two_data",
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := syncAllVideoToElasticsearchLogic.SyncAllVideoToElasticsearch()
			if err != nil {
				assert.Equal(t, err.Error(), tc.err.Error())
			} else {
				assert.Equal(t, err, tc.err)
			}
		})
	}
}
