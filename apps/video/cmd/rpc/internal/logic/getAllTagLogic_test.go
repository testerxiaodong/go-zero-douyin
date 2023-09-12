package logic_test

import (
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/syncx"
	"go-zero-douyin/apps/video/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/mock"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestGetAllTagLogic_GetAllTag(t *testing.T) {
	ctl := gomock.NewController(t)

	defer ctl.Finish()

	mockTagDo := mock.NewMockTagDo(ctl)

	mockRedis := globalMock.NewMockRedisCache(ctl)

	mockRabbit := globalMock.NewMockSender(ctl)

	utils.IgnoreGo()
	defer utils.RecoverGo()

	serviceContext := &svc.ServiceContext{TagDo: mockTagDo, Redis: mockRedis,
		Rabbit: mockRabbit, SingleFlight: syncx.NewSingleFlight()}

	getAllTagLogic := logic.NewGetAllTagLogic(context.Background(), serviceContext)

	// redis中有数据的mock
	expectedValue := NewRandTag()
	tagString, err := json.Marshal(expectedValue)
	require.NoError(t, err)
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
	mockRedis.EXPECT().Smembers(gomock.Any(), gomock.Any()).Return([]string{string(tagString), string(tagString)}, nil)
	mockRedis.EXPECT().Expire(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// redis中没有数据，查询数据库失败mock
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
	dbSearchError := errors.New("TagDo.GetAllTags error")
	mockTagDo.EXPECT().GetAllTags(gomock.Any()).Return(nil, dbSearchError)

	// redis没有数据，查询数据库成功的，但数据库没有数据的mock
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
	mockTagDo.EXPECT().GetAllTags(gomock.Any()).
		Return([]*model.Tag{}, nil)

	// redis没有数据，查询数据库成功的，且数据库有数据的mock
	mockRedis.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
	mockTagDo.EXPECT().GetAllTags(gomock.Any()).
		Return([]*model.Tag{expectedValue, expectedValue}, nil)

	// 表格驱动测试
	testcases := []struct {
		name string
		req  *pb.GetAllTagReq
		err  error
	}{
		{
			name: "get_all_tag_with_redis",
			req:  &pb.GetAllTagReq{},
			err:  nil,
		},
		{
			name: "get_all_tag_with_database_search_error",
			req:  &pb.GetAllTagReq{},
			err:  errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "从数据库中获取所有视频标签信息失败, err: %v", dbSearchError),
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
		ID:   utils.NewRandomInt64(1, 10),
		Name: utils.NewRandomString(10),
	}
}
