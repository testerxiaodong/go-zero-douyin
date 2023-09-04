package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/video/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/mock"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"testing"
)

func TestPublishVideoLogic_PublishVideo(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockVideoDo := mock.NewMockVideoDo(ctl)

	serviceContext := &svc.ServiceContext{VideoDo: mockVideoDo}

	publishVideoLogic := logic.NewPublishVideoLogic(context.Background(), serviceContext)

	expectVideo := NewRandomVideo()

	// 插入数据库失败mock
	mockVideoDo.EXPECT().InsertVideo(gomock.Any(), gomock.Any()).Return(errors.New("insert database error"))

	// 插入数据库成功
	mockVideoDo.EXPECT().InsertVideo(gomock.Any(), gomock.Any()).Return(nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.PublishVideoReq
	}{
		{
			name: "publish_video_with_empty_param",
			req:  nil,
		},
		{
			name: "publish_video_with_insert_database_error",
			req:  &pb.PublishVideoReq{Title: expectVideo.Title, OwnerId: expectVideo.OwnerID, PlayUrl: expectVideo.PlayURL, CoverUrl: expectVideo.CoverURL},
		},
		{
			name: "publish_video_with_insert_database_error",
			req:  &pb.PublishVideoReq{Title: expectVideo.Title, OwnerId: expectVideo.OwnerID, PlayUrl: expectVideo.PlayURL, CoverUrl: expectVideo.CoverURL},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := publishVideoLogic.PublishVideo(tc.req)
			if err == nil {
				assert.NotEmpty(t, resp)
				assert.Equal(t, resp.GetVideo().GetTitle(), expectVideo.Title)
				assert.Equal(t, resp.GetVideo().GetOwnerId(), expectVideo.OwnerID)
				assert.Equal(t, resp.GetVideo().GetPlayUrl(), expectVideo.PlayURL)
				assert.Equal(t, resp.GetVideo().GetCoverUrl(), expectVideo.CoverURL)
			} else {
				assert.Empty(t, resp)
			}
		})
	}
}
