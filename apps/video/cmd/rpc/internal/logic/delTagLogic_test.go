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

func TestDelTagLogic_DelTag(t *testing.T) {
	ctl := gomock.NewController(t)

	mockTagDo := mock.NewMockTagDo(ctl)

	mockRedis := globalMock.NewMockRedisCache(ctl)

	mockRabbit := globalMock.NewMockSender(ctl)

	serviceContext := &svc.ServiceContext{TagDo: mockTagDo, Redis: mockRedis, Rabbit: mockRabbit}

	delTagLogic := logic.NewDelTagLogic(context.Background(), serviceContext)

	// 删除失败的mock
	dbDeleteError := errors.New("TagDo.DeleteTag error")
	mockTagDo.EXPECT().DeleteTag(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, dbDeleteError)

	// 删除缓存失败，且发送消息失败的mock
	redisDeleteError := errors.New("redis delete error")
	senderError := errors.New("send message error")
	mockTagDo.EXPECT().DeleteTag(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisDeleteError)
	mockRabbit.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 删除缓存失败，但发送消息成功的mock
	mockTagDo.EXPECT().DeleteTag(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisDeleteError)
	mockRabbit.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 删除缓存成功的mock
	mockTagDo.EXPECT().DeleteTag(gomock.Any(), gomock.Any()).Return(gen.ResultInfo{}, nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.DelTagReq
		err  error
	}{
		{
			name: "add_section_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del tag with empty param"),
		},
		{
			name: "add_section_with_empty_name",
			req:  &pb.DelTagReq{Id: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del tag with empty id"),
		},
		{
			name: "add_section_with_database_delete_error",
			req:  &pb.DelTagReq{Id: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "删除标签失败, err: %v, id: %d", dbDeleteError, 1),
		},
		{
			name: "add_section_with_sender_error",
			req:  &pb.DelTagReq{Id: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("发布删除标签缓存信息失败"), "err: %v", senderError),
		},
		{
			name: "add_section_with_redis_delete_error",
			req:  &pb.DelTagReq{Id: 1},
			err:  nil,
		},
		{
			name: "add_section_success",
			req:  &pb.DelTagReq{Id: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := delTagLogic.DelTag(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
