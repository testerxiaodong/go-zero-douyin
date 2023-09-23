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

func TestNewGetTagByIdLogic(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockTagDo := mock.NewMocktagModel(ctl)
	serviceContext := &svc.ServiceContext{TagModel: mockTagDo}
	getTagByIdLogic := logic.NewGetTagByIdLogic(context.Background(), serviceContext)

	// 查询数据库失败的mock
	dbSearchError := errors.New("TagDo.GetTagById error")
	mockTagDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, dbSearchError)

	// 查询数据库成功但没有数据的mock
	mockTagDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(nil, model.ErrNotFound)

	// 查询数据库成功且有数据的mock
	expectedValue := NewRandTag()
	mockTagDo.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(expectedValue, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *pb.GetTagByIdReq
		err  error
	}{
		{
			name: "get_tag_by_id_with_empty_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get tag by id with empty param"),
		},
		{
			name: "get_tag_by_id_with_empty_id",
			req:  &pb.GetTagByIdReq{Id: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get tag by id with empty id"),
		},
		{
			name: "get_tag_by_id_with_database_search_error",
			req:  &pb.GetTagByIdReq{Id: 1},
			err: errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"通过id从数据库中获取标签信息失败, err: %v, id: %d", dbSearchError, 1),
		},
		{
			name: "get_tag_by_id_with_no_database_record",
			req:  &pb.GetTagByIdReq{Id: 1},
			err:  errors.Wrapf(xerr.NewErrMsg("标签不存在"), "tag_id: %d", 1),
		},
		{
			name: "get_tag_by_id_success",
			req:  &pb.GetTagByIdReq{Id: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := getTagByIdLogic.GetTagById(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
				assert.Equal(t, resp.GetTag().GetId(), expectedValue.Id)
				assert.Equal(t, resp.GetTag().GetName(), expectedValue.Name)
			}
		})
	}
}
