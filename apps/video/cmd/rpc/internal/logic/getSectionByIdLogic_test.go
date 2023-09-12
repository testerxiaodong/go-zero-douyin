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
	"gorm.io/gorm"
	"testing"
)

func TestGetSectionByIdLogic_GetSectionById(t *testing.T) {
	ctl := gomock.NewController(t)

	mockSectionDo := mock.NewMockSectionDo(ctl)

	serviceContext := &svc.ServiceContext{SectionDo: mockSectionDo}

	getSectionByIdLogic := logic.NewGetSectionByIdLogic(context.Background(), serviceContext)

	// 查询数据库失败的mock
	dbSearchError := errors.New("SectionDo.GetSectionById error")
	mockSectionDo.EXPECT().GetSectionById(gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 查询数据库成功但没有数据的mock
	mockSectionDo.EXPECT().GetSectionById(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

	// 查询数据库成功且有数据的mock
	expectedValue := NewRandSection()
	mockSectionDo.EXPECT().GetSectionById(gomock.Any(), gomock.Any()).Return(expectedValue, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.GetSectionByIdReq
		err  error
	}{
		{
			name: "get_section_by_id_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get section by id with empty param"),
		},
		{
			name: "get_section_by_id_with_empty_id",
			req:  &pb.GetSectionByIdReq{Id: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get section by id with empty id"),
		},
		{
			name: "get_section_by_id_with_database_search_error",
			req:  &pb.GetSectionByIdReq{Id: 1},
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"通过id从数据库中获取分区信息失败, err: %v, id: %d", dbSearchError, 1),
		},
		{
			name: "get_section_by_id_with_no_database_record",
			req:  &pb.GetSectionByIdReq{Id: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("分区不存在"), "section_id: %d", 1),
		},
		{
			name: "get_section_by_id_success",
			req:  &pb.GetSectionByIdReq{Id: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := getSectionByIdLogic.GetSectionById(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
				assert.Equal(t, resp.GetSectionInfo().GetId(), expectedValue.ID)
				assert.Equal(t, resp.GetSectionInfo().GetName(), expectedValue.Name)
			}
		})
	}
}
