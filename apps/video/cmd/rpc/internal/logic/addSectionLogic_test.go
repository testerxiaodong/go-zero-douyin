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
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"gorm.io/gorm"
	"testing"
)

func TestNewAddSectionLogic(t *testing.T) {
	ctl := gomock.NewController(t)

	mockSectionDo := mock.NewMockSectionDo(ctl)

	mockRedis := globalMock.NewMockRedisCache(ctl)

	mockRabbit := globalMock.NewMockSender(ctl)

	serviceContext := &svc.ServiceContext{SectionDo: mockSectionDo, Redis: mockRedis, Rabbit: mockRabbit}

	addSectionLogic := logic.NewAddSectionLogic(context.Background(), serviceContext)

	// 查询分区失败的mock
	dbSearchError := errors.New("SectionDo.GetSectionByName error")
	mockSectionDo.EXPECT().GetSectionByName(gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 分区已存在mock
	mockSectionDo.EXPECT().GetSectionByName(gomock.Any(), gomock.Any()).Return(&model.Section{ID: 1, Name: "test"}, nil)

	// 插入失败的mock
	dbInsertError := errors.New("SectionDo.InsertSection error")
	mockSectionDo.EXPECT().GetSectionByName(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockSectionDo.EXPECT().InsertSection(gomock.Any(), gomock.Any()).Return(dbInsertError)

	// 删除缓存失败，且发送消息失败的mock
	redisDeleteError := errors.New("redis delete error")
	senderError := errors.New("send message error")
	mockSectionDo.EXPECT().GetSectionByName(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockSectionDo.EXPECT().InsertSection(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisDeleteError)
	mockRabbit.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 删除缓存失败，但发送消息成功的mock
	mockSectionDo.EXPECT().GetSectionByName(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockSectionDo.EXPECT().InsertSection(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisDeleteError)
	mockRabbit.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 删除缓存成功的mock
	mockSectionDo.EXPECT().GetSectionByName(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockSectionDo.EXPECT().InsertSection(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.AddSectionReq
		err  error
	}{
		{
			name: "add_section_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "add section with empty param"),
		},
		{
			name: "add_section_with_empty_name",
			req:  &pb.AddSectionReq{Name: ""},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "add section with empty name"),
		},
		{
			name: "add_section_with_database_search_error",
			req:  &pb.AddSectionReq{Name: "test"},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "数据库根据名称查询分区失败, err: %v, name: %s", dbSearchError, "test"),
		},
		{
			name: "add_section_with_exist_record",
			req:  &pb.AddSectionReq{Name: "test"},
			err:  errors.Wrapf(xerr.NewErrMsg("分区名已存在"), "新增分区失败， req: %v", &pb.AddSectionReq{Name: "test"}),
		},
		{
			name: "add_section_with_database_insert_error",
			req:  &pb.AddSectionReq{Name: "test"},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "数据库新增分区失败, err: %v, name: %s", dbInsertError, "test"),
		},
		{
			name: "add_section_with_sender_error",
			req:  &pb.AddSectionReq{Name: "test"},
			err:  errors.Wrapf(xerr.NewErrMsg("发布删除分区缓存信息失败"), "err: %v", senderError),
		},
		{
			name: "add_section_with_redis_delete_error",
			req:  &pb.AddSectionReq{Name: "test"},
			err:  nil,
		},
		{
			name: "add_section_success",
			req:  &pb.AddSectionReq{Name: "test"},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := addSectionLogic.AddSection(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
