package logic_test

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/video/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/mock"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"testing"
)

func TestGetAllSectionLogic_GetAllSection(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockSectionDo := mock.NewMocksectionModel(ctl)
	serviceContext := &svc.ServiceContext{SectionModel: mockSectionDo}
	getAllSectionLogic := logic.NewGetAllSectionLogic(context.Background(), serviceContext)

	// 查询失败的mock
	dbSearchError := errors.New("db search error")
	mockSectionDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockSectionDo.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 插叙成功长度为零的mock
	mockSectionDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockSectionDo.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*model.Section{}, nil)

	// 查询成功，长度为2的mock
	mockSectionDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockSectionDo.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*model.Section{NewRandSection(), NewRandSection()}, nil)
	// 表格驱动测试
	testcases := []struct {
		name string
		req  *pb.GetAllSectionReq
		err  error
	}{
		{
			name: "get_all_section_with_database_search_error",
			req:  &pb.GetAllSectionReq{},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "查询所有视频分区失败, err: %v", dbSearchError),
		},
		{
			name: "get_all_section_with_database_no_record",
			req:  &pb.GetAllSectionReq{},
			err:  nil,
		},
		{
			name: "get_all_section_with_database_two_record",
			req:  &pb.GetAllSectionReq{},
			err:  nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getAllSectionLogic.GetAllSection(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}

func NewRandSection() *model.Section {
	return &model.Section{
		Id:   utils.NewRandomInt64(1, 10),
		Name: utils.NewRandomString(10),
	}
}
