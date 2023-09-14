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
	"strings"
	"testing"
)

func TestPublishVideoLogic_PublishVideo(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockVideoDo := mock.NewMockVideoDo(ctl)
	mockSender := globalMock.NewMockSender(ctl)
	serviceContext := &svc.ServiceContext{VideoDo: mockVideoDo, Rabbit: mockSender}

	publishVideoLogic := logic.NewPublishVideoLogic(context.Background(), serviceContext)

	expectVideo := NewRandomVideo()

	// 插入数据库失败mock
	dbError := errors.New("insert database error")
	mockVideoDo.EXPECT().InsertVideo(gomock.Any(), gomock.Any()).Return(dbError)

	// 插入数据库成功，发送消息失败的mock
	senderError := errors.New("sender error")
	mockVideoDo.EXPECT().InsertVideo(gomock.Any(), gomock.Any()).Return(nil)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 成功的mock
	mockVideoDo.EXPECT().InsertVideo(gomock.Any(), gomock.Any()).Return(nil)
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 表格驱动测试
	req := &pb.PublishVideoReq{Title: expectVideo.Title, OwnerId: expectVideo.OwnerID, SectionId: expectVideo.SectionID,
		Tags: strings.Split(expectVideo.TagIds, ","), PlayUrl: expectVideo.PlayURL, CoverUrl: expectVideo.CoverURL}
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
			name: "publish_video_with_insert_database_error",
			req:  req,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "insert video failed, err: %v", dbError),
		},
		{
			name: "publish_video_with_sender_error",
			req:  req,
			err:  errors.Wrapf(xerr.NewErrCode(xerr.RPC_UPDATE_ERR), "req: %v, err: %v", req, senderError),
		},
		{
			name: "publish_video_with_success",
			req:  req,
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
