package logic_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/recommend/cmd/rpc/internal/logic"
	"go-zero-douyin/apps/recommend/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/recommend/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"go-zero-douyin/mock"
	"testing"
)

func TestNewVideoRecommendSectionLogic(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockRest := mock.NewMockRestClient(ctl)
	serviceContext := &svc.ServiceContext{RestClient: mockRest}
	VideoRecommendSectionLogic := logic.NewVideoRecommendSectionLogic(context.Background(), serviceContext)

	// GorseApi调用失败的mock
	gorseError := errors.New("gorse error")
	mockRest.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(gorseError)

	// GorseApi调用成功的mock
	mockRest.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 表格驱动测试
	req := &pb.VideoRecommendSectionReq{
		UserId:    utils.NewRandomInt64(1, 10),
		SectionId: utils.NewRandomInt64(1, 10),
		Count:     utils.NewRandomInt64(1, 10),
	}
	testCases := []struct {
		name string
		req  *pb.VideoRecommendSectionReq
		err  error
	}{
		{
			name: "video_recommend_section_with_nil_param",
			req:  nil,
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "获取分区视频推荐时参数为nil"),
		},
		{
			name: "video_recommend_section_with_empty_param",
			req:  &pb.VideoRecommendSectionReq{UserId: 0, SectionId: 0, Count: 0},
			err:  errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "获取分区视频推荐时分区id或者用户id不能为空"),
		},
		{
			name: "video_recommend_section_with_gorse_error",
			req:  req,
			err: errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR),
				"GorseApi获取推荐失败, err: %v, req: %v", gorseError, req),
		},
		{
			name: "video_recommend_section_success",
			req:  req,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := VideoRecommendSectionLogic.VideoRecommendSection(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
