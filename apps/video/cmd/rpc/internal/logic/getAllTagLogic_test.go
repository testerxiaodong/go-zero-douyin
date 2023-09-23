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

func TestGetAllTagLogic_GetAllTag(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockTagDo := mock.NewMocktagModel(ctl)
	serviceContext := &svc.ServiceContext{TagModel: mockTagDo}
	getAllTagLogic := logic.NewGetAllTagLogic(context.Background(), serviceContext)

	// 查询失败的mock
	dbSearchError := errors.New("db search error")
	mockTagDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockTagDo.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 插叙成功长度为零的mock
	mockTagDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockTagDo.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*model.Tag{}, nil)

	// 查询成功，长度为2的mock
	mockTagDo.EXPECT().SelectBuilder().Return(squirrel.SelectBuilder{})
	mockTagDo.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*model.Tag{NewRandTag(), NewRandTag()}, nil)
	// 表格驱动测试
	testcases := []struct {
		name string
		req  *pb.GetAllTagReq
		err  error
	}{
		{
			name: "get_all_tag_with_database_search_error",
			req:  &pb.GetAllTagReq{},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "查询所有视频标签失败, err: %v", dbSearchError),
		},
		{
			name: "get_all_tag_with_database_no_record",
			req:  &pb.GetAllTagReq{},
			err:  nil,
		},
		{
			name: "get_all_tag_with_database_two_record",
			req:  &pb.GetAllTagReq{},
			err:  nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getAllTagLogic.GetAllTag(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}

func NewRandTag() *model.Tag {
	return &model.Tag{
		Id:   utils.NewRandomInt64(1, 10),
		Name: utils.NewRandomString(10),
	}
}
