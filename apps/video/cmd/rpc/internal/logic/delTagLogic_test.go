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

func TestDelTagLogic_DelTag(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockTagDo := mock.NewMocktagModel(ctl)
	serviceContext := &svc.ServiceContext{TagModel: mockTagDo}
	delTagLogic := logic.NewDelTagLogic(context.Background(), serviceContext)

	// 查询失败的mock
	dbSearchError := errors.New("TagDo.FindOne error")
	dbDeleteError := errors.New("TagDo.DeleteSoft error")
	mockTagDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 标签不存在的mock
	mockTagDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, model.ErrNotFound)

	// 删除失败的mock
	mockTagDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(&model.Tag{}, nil)
	mockTagDo.EXPECT().DeleteSoft(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbDeleteError)

	// 删除成功的mock
	mockTagDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(&model.Tag{}, nil)
	mockTagDo.EXPECT().DeleteSoft(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.DelTagReq
		err  error
	}{
		{
			name: "add_tag_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del tag with empty param"),
		},
		{
			name: "add_tag_with_empty_name",
			req:  &pb.DelTagReq{Id: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del tag with empty id"),
		},
		{
			name: "add_tag_with_database_search_error",
			req:  &pb.DelTagReq{Id: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "查询视频标签失败, tag_id: %d", 1),
		},
		{
			name: "add_tag_with_no_tag_error",
			req:  &pb.DelTagReq{Id: 1},
			err:  errors.Wrapf(logic.ErrTagNotFound, "tag_id: %d", 1),
		},
		{
			name: "add_tag_with_db_delete_error",
			req:  &pb.DelTagReq{Id: 1},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "删除标签失败, err: %v, id: %d", dbDeleteError, 1),
		},
		{
			name: "add_tag_success",
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
