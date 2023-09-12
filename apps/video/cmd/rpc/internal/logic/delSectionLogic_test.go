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
	"gorm.io/gen"
	"testing"
)

func TestDelSectionLogic_DelSection(t *testing.T) {
	ctl := gomock.NewController(t)

	mockSectionDo := mock.NewMockSectionDo(ctl)

	mockRedis := globalMock.NewMockRedisCache(ctl)

	mockRabbit := globalMock.NewMockSender(ctl)

	serviceContext := &svc.ServiceContext{SectionDo: mockSectionDo, Redis: mockRedis, Rabbit: mockRabbit}

	delSectionLogic := logic.NewDelSectionLogic(context.Background(), serviceContext)

	// 删除失败的mock
	dbDeleteError := errors.New("SectionDo.DeleteSection error")
	mockSectionDo.EXPECT().DeleteSection(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, dbDeleteError)

	// 删除缓存失败，且发送消息失败的mock
	redisDeleteError := errors.New("redis delete error")
	senderError := errors.New("send message error")
	mockSectionDo.EXPECT().DeleteSection(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisDeleteError)
	mockRabbit.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 删除缓存失败，但发送消息成功的mock
	mockSectionDo.EXPECT().DeleteSection(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisDeleteError)
	mockRabbit.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 删除缓存成功的mock
	mockSectionDo.EXPECT().DeleteSection(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.DelSectionReq
		err  error
	}{
		{
			name: "add_section_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del section with empty param"),
		},
		{
			name: "add_section_with_empty_name",
			req:  &pb.DelSectionReq{Id: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del section with empty id"),
		},
		{
			name: "add_section_with_database_delete_error",
			req:  &pb.DelSectionReq{Id: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "删除分区失败, err: %v, id: %d", dbDeleteError, 1),
		},
		{
			name: "add_section_with_sender_error",
			req:  &pb.DelSectionReq{Id: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("发布删除分区缓存信息失败"), "err: %v", senderError),
		},
		{
			name: "add_section_with_redis_delete_error",
			req:  &pb.DelSectionReq{Id: 1},
			err:  nil,
		},
		{
			name: "add_section_success",
			req:  &pb.DelSectionReq{Id: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := delSectionLogic.DelSection(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
