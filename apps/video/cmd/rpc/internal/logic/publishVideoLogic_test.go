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
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestPublishVideoLogic_PublishVideo(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockVideoDo := mock.NewMockvideoModel(ctl)
	mockAsynq := globalMock.NewMockTaskQueueClient(ctl)
	serviceContext := &svc.ServiceContext{VideoModel: mockVideoDo, Asynq: mockAsynq}
	publishVideoLogic := logic.NewPublishVideoLogic(context.Background(), serviceContext)
	expectVideo := NewRandomVideo()

	// 延迟任务发布失败mock
	asynqError := errors.New("asynq enqueue error")
	mockAsynq.EXPECT().EnqueueContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, asynqError)

	// 延迟任务发布成功mock
	mockAsynq.EXPECT().EnqueueContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

	// 插入数据库失败mock
	dbError := errors.New("insert database error")
	result := globalMock.NewMockResult(ctl)
	mockVideoDo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(result, dbError)

	// 成功的mock
	mockVideoDo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(result, nil)

	// 表格驱动测试
	reqWithPublishTime := &pb.PublishVideoReq{Title: expectVideo.Title, OwnerId: expectVideo.OwnerId, SectionId: expectVideo.SectionId,
		TagIds: expectVideo.TagIds, PlayUrl: expectVideo.PlayUrl, CoverUrl: expectVideo.CoverUrl, PublishTime: 1}
	reqWithoutPublishTime := &pb.PublishVideoReq{Title: expectVideo.Title, OwnerId: expectVideo.OwnerId, SectionId: expectVideo.SectionId,
		TagIds: expectVideo.TagIds, PlayUrl: expectVideo.PlayUrl, CoverUrl: expectVideo.CoverUrl}
	testCases := []struct {
		name string
		req  *pb.PublishVideoReq
		err  error
	}{
		{
			name: "publish_video_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Publish video empty param"),
		},
		{
			name: "publish_video_with_asynq_error",
			req:  reqWithPublishTime,
			err:  errors.Wrapf(xerr.NewErrMsg("创建发布视频的延迟任务失败"), "err: %v", asynqError),
		},
		{
			name: "publish_video_with_asynq_success",
			req:  reqWithPublishTime,
			err:  nil,
		},
		{
			name: "publish_video_with_insert_database_error",
			req:  reqWithoutPublishTime,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "insert video failed, err: %v", dbError),
		},
		{
			name: "publish_video_with_success",
			req:  reqWithoutPublishTime,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := publishVideoLogic.PublishVideo(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
