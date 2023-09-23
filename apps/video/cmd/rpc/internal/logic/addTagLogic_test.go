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
	"testing"
)

func TestAddTagLogic_AddTag(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockTagDo := mock.NewMocktagModel(ctl)
	serviceContext := &svc.ServiceContext{TagModel: mockTagDo}
	addTagLogic := logic.NewAddTagLogic(context.Background(), serviceContext)

	// 查询标签失败的mock
	dbSearchError := errors.New("TagDo.GetTagByName error")
	mockTagDo.EXPECT().FindOneByName(gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 标签已存在mock
	mockTagDo.EXPECT().FindOneByName(gomock.Any(), gomock.Any()).Return(&model.Tag{Id: 1, Name: "test"}, nil)

	// 插入失败的mock
	dbInsertError := errors.New("TagDo.InsertTag error")
	mockTagDo.EXPECT().FindOneByName(gomock.Any(), gomock.Any()).Return(nil, model.ErrNotFound)
	mockTagDo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbInsertError)

	// 插入成功的mock
	mockTagDo.EXPECT().FindOneByName(gomock.Any(), gomock.Any()).Return(nil, model.ErrNotFound)
	mockTagDo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

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
