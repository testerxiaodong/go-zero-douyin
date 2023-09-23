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
	"go-zero-douyin/common/xerr"
	"strings"
	"testing"
)

func TestGetVideoByIdLogic_GetVideoById(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockVideoDo := mock.NewMockvideoModel(ctl)

	serviceContext := &svc.ServiceContext{VideoModel: mockVideoDo}

	getVideoByIdLogic := logic.NewGetVideoByIdLogic(context.Background(), serviceContext)

	mockVideo := NewRandomVideo()

	// 查询数据库失败mock
	dbError := errors.New("search database error")
	mockVideoDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, dbError)

	// 视频不存在mock
	mockVideoDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, model.ErrNotFound)

	// 查询视频成功mock
	mockVideoDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(mockVideo, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.GetVideoByIdReq
		err  error
	}{
		{
			name: "get_video_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get video by id with empty param"),
		},
		{
			name: "get_video_with_empty_id",
			req:  &pb.GetVideoByIdReq{Id: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get video by id with empty id"),
		},
		{
			name: "get_video_with_search_database_error",
			req:  &pb.GetVideoByIdReq{Id: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "find video by id failed: %v", dbError),
		},
		{
			name: "get_video_with_no_database_record",
			req:  &pb.GetVideoByIdReq{Id: 1},
			err:  errors.Wrapf(logic.ErrVideoNotFound, "video_id: %d", 1),
		},
		{
			name: "get_video_success",
			req:  &pb.GetVideoByIdReq{Id: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			infoResp, err := getVideoByIdLogic.GetVideoById(tc.req)
			if err == nil {
				assert.NotEmpty(t, infoResp)
				assert.Equal(t, infoResp.GetVideo().GetId(), mockVideo.Id)
				assert.Equal(t, infoResp.GetVideo().GetTitle(), mockVideo.Title)
				assert.Equal(t, infoResp.GetVideo().GetSectionId(), mockVideo.SectionId)
				assert.Equal(t, infoResp.GetVideo().GetTagIds(), mockVideo.TagIds)
				assert.Equal(t, infoResp.GetVideo().GetOwnerId(), mockVideo.OwnerId)
				assert.Equal(t, infoResp.GetVideo().GetOwnerName(), mockVideo.OwnerName)
				assert.Equal(t, infoResp.GetVideo().GetPlayUrl(), mockVideo.PlayUrl)
				assert.Equal(t, infoResp.GetVideo().GetCoverUrl(), mockVideo.CoverUrl)
			} else {
				assert.Equal(t, tc.err.Error(), err.Error())
			}
		})
	}
}

func NewRandomVideo() *model.Video {
	video := &model.Video{}
	video.Id = utils.NewRandomInt64(1, 10)
	video.Title = utils.NewRandomString(10)
	video.SectionId = utils.NewRandomInt64(1, 10)
	video.TagIds = strings.Join([]string{"11", "22"}, ",")
	video.OwnerId = utils.NewRandomInt64(1, 10)
	video.OwnerName = utils.NewRandomString(10)
	video.PlayUrl = utils.NewRandomString(10)
	video.CoverUrl = utils.NewRandomString(10)

	return video
}
