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

func TestDelSectionLogic_DelSection(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockSectionDo := mock.NewMocksectionModel(ctl)
	serviceContext := &svc.ServiceContext{SectionModel: mockSectionDo}
	delSectionLogic := logic.NewDelSectionLogic(context.Background(), serviceContext)

	// 查询失败的mock
	dbSearchError := errors.New("SectionDo.FindOne error")
	dbDeleteError := errors.New("SectionDo.DeleteSoft error")
	mockSectionDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 分区不存在的mock
	mockSectionDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, model.ErrNotFound)

	// 删除失败的mock
	mockSectionDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(&model.Section{}, nil)
	mockSectionDo.EXPECT().DeleteSoft(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbDeleteError)

	// 删除成功的mock
	mockSectionDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(&model.Section{}, nil)
	mockSectionDo.EXPECT().DeleteSoft(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

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
			name: "add_section_with_database_search_error",
			req:  &pb.DelSectionReq{Id: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "查询视频分区失败, section_id: %d", 1),
		},
		{
			name: "add_section_with_no_section_error",
			req:  &pb.DelSectionReq{Id: 1},
			err:  errors.Wrapf(logic.ErrSectionNotFound, "分区不存在, section_id: %d", 1),
		},
		{
			name: "add_section_with_db_delete_error",
			req:  &pb.DelSectionReq{Id: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "删除分区失败, err: %v, id: %d", dbDeleteError, 1),
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
