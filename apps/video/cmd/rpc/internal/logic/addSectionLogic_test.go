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

func TestNewAddSectionLogic(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockSectionDo := mock.NewMocksectionModel(ctl)
	serviceContext := &svc.ServiceContext{SectionModel: mockSectionDo}
	addSectionLogic := logic.NewAddSectionLogic(context.Background(), serviceContext)

	// 查询分区失败的mock
	dbSearchError := errors.New("SectionDo.GetSectionByName error")
	mockSectionDo.EXPECT().FindOneByNameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 分区已存在mock
	mockSectionDo.EXPECT().FindOneByNameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Section{Id: 1, Name: "test"}, nil)

	// 插入失败的mock
	dbInsertError := errors.New("SectionDo.InsertSection error")
	mockSectionDo.EXPECT().FindOneByNameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, model.ErrNotFound)
	mockSectionDo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbInsertError)

	// 插入成功的mock
	mockSectionDo.EXPECT().FindOneByNameIsDelete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, model.ErrNotFound)
	mockSectionDo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

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
