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

func TestAddTagLogic_AddTag(t *testing.T) {
	ctl := gomock.NewController(t)

	mockTagDo := mock.NewMockTagDo(ctl)

	mockRedis := globalMock.NewMockRedisCache(ctl)

	mockRabbit := globalMock.NewMockSender(ctl)

	serviceContext := &svc.ServiceContext{TagDo: mockTagDo, Redis: mockRedis, Rabbit: mockRabbit}

	addTagLogic := logic.NewAddTagLogic(context.Background(), serviceContext)

	// 查询分区失败的mock
	dbSearchError := errors.New("TagDo.GetTagByName error")
	mockTagDo.EXPECT().GetTagByName(gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 分区已存在mock
	mockTagDo.EXPECT().GetTagByName(gomock.Any(), gomock.Any()).Return(&model.Tag{ID: 1, Name: "test"}, nil)

	// 插入失败的mock
	dbInsertError := errors.New("TagDo.InsertTag error")
	mockTagDo.EXPECT().GetTagByName(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockTagDo.EXPECT().InsertTag(gomock.Any(), gomock.Any()).Return(dbInsertError)

	// 删除缓存失败，且发送消息失败的mock
	redisDeleteError := errors.New("redis delete error")
	senderError := errors.New("send message error")
	mockTagDo.EXPECT().GetTagByName(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockTagDo.EXPECT().InsertTag(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisDeleteError)
	mockRabbit.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(senderError)

	// 删除缓存失败，但发送消息成功的mock
	mockTagDo.EXPECT().GetTagByName(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockTagDo.EXPECT().InsertTag(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(0, redisDeleteError)
	mockRabbit.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 删除缓存成功的mock
	mockTagDo.EXPECT().GetTagByName(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	mockTagDo.EXPECT().InsertTag(gomock.Any(), gomock.Any()).Return(nil)
	mockRedis.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(1, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.AddTagReq
		err  error
	}{
		{
			name: "add_tag_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "add tag with empty param"),
		},
		{
			name: "add_tag_with_empty_name",
			req:  &pb.AddTagReq{Name: ""},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "add tag with empty name"),
		},
		{
			name: "add_tag_with_database_search_error",
			req:  &pb.AddTagReq{Name: "test"},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "数据库根据名称查询标签失败, err: %v, name: %s", dbSearchError, "test"),
		},
		{
			name: "add_tag_with_exist_record",
			req:  &pb.AddTagReq{Name: "test"},
			err:  errors.Wrapf(xerr.NewErrMsg("标签名已存在"), "新增标签失败， req: %v", &pb.AddTagReq{Name: "test"}),
		},
		{
			name: "add_tag_with_database_insert_error",
			req:  &pb.AddTagReq{Name: "test"},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "数据库新增标签失败, err: %v, name: %s", dbInsertError, "test"),
		},
		{
			name: "add_tag_with_sender_error",
			req:  &pb.AddTagReq{Name: "test"},
			err:  errors.Wrapf(xerr.NewErrMsg("发布删除标签缓存信息失败"), "err: %v", senderError),
		},
		{
			name: "add_tag_with_redis_delete_error",
			req:  &pb.AddTagReq{Name: "test"},
			err:  nil,
		},
		{
			name: "add_tag_success",
			req:  &pb.AddTagReq{Name: "test"},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := addTagLogic.AddTag(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
