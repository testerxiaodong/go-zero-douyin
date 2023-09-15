package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/video/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/mock"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"gorm.io/gorm"
	"strings"
	"testing"
)

func TestGetVideoByIdLogic_GetVideoById(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockVideoDo := mock.NewMockVideoDo(ctl)

	serviceContext := &svc.ServiceContext{VideoDo: mockVideoDo}

	getVideoByIdLogic := logic.NewGetVideoByIdLogic(context.Background(), serviceContext)

	mockVideo := NewRandomVideo()

	// 查询数据库失败mock
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(nil, errors.New("search database error"))

	// 视频不存在mock
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	// 查询视频成功mock
	mockVideoDo.EXPECT().GetVideoById(gomock.Any(), gomock.Any()).Return(mockVideo, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.GetVideoByIdReq
	}{
		{
			name: "get_video_with_empty_param",
			req:  nil,
		},
		{
			name: "get_video_with_empty_id",
			req:  &pb.GetVideoByIdReq{Id: 0},
		},
		{
			name: "get_video_with_search_database_error",
			req:  &pb.GetVideoByIdReq{Id: 1},
		},
		{
			name: "get_video_with_no_database_record",
			req:  &pb.GetVideoByIdReq{Id: 1},
		},
		{
			name: "get_video_success",
			req:  &pb.GetVideoByIdReq{Id: 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			infoResp, err := getVideoByIdLogic.GetVideoById(tc.req)
			if err == nil {
				assert.NotEmpty(t, infoResp)
				assert.Equal(t, infoResp.GetVideo().GetId(), mockVideo.ID)
				assert.Equal(t, infoResp.GetVideo().GetTitle(), mockVideo.Title)
				assert.Equal(t, infoResp.GetVideo().GetSectionId(), mockVideo.SectionID)
				assert.Equal(t, infoResp.GetVideo().GetTags(), strings.Split(mockVideo.TagIds, ","))
				assert.Equal(t, infoResp.GetVideo().GetOwnerId(), mockVideo.OwnerID)
				assert.Equal(t, infoResp.GetVideo().GetOwnerName(), mockVideo.OwnerName)
				assert.Equal(t, infoResp.GetVideo().GetPlayUrl(), mockVideo.PlayURL)
				assert.Equal(t, infoResp.GetVideo().GetCoverUrl(), mockVideo.CoverURL)
			} else {
				assert.Empty(t, infoResp)
			}
		})
	}
}

func NewRandomVideo() *model.Video {
	video := &model.Video{}
	video.ID = utils.NewRandomInt64(1, 10)
	video.Title = utils.NewRandomString(10)
	video.SectionID = utils.NewRandomInt64(1, 10)
	video.TagIds = strings.Join([]string{"11", "22"}, ",")
	video.OwnerID = utils.NewRandomInt64(1, 10)
	video.OwnerName = utils.NewRandomString(10)
	video.PlayURL = utils.NewRandomString(10)
	video.CoverURL = utils.NewRandomString(10)

	return video
}
