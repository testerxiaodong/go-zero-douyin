package search_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/search/cmd/api/internal/logic/search"
	"go-zero-douyin/apps/search/cmd/api/internal/svc"
	"go-zero-douyin/apps/search/cmd/api/internal/types"
	searchMock "go-zero-douyin/apps/search/cmd/rpc/mock"
	"go-zero-douyin/apps/search/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	globalMock "go-zero-douyin/mock"
	"testing"
)

func TestVideoLogic_Video(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockValidator := globalMock.NewMockValidator(ctl)
	mockSearchRpc := searchMock.NewMockSearch(ctl)
	serviceContext := &svc.ServiceContext{Validator: mockValidator, SearchRpc: mockSearchRpc}
	videoLogic := search.NewVideoLogic(context.Background(), serviceContext)

	// 参数校验失败的mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// searchRpc.SearchVideo失败的mock
	esError := errors.New("es error")
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockSearchRpc.EXPECT().SearchVideo(gomock.Any(), gomock.Any()).Return(nil, esError)

	// 成功的mock
	resp := &pb.SearchVideoResp{Total: 1, Videos: []*pb.Video{&pb.Video{
		Id:           utils.NewRandomInt64(1, 10),
		Title:        utils.NewRandomString(10),
		SectionId:    utils.NewRandomInt64(1, 10),
		Tags:         []string{"1", "2", "3"},
		OwnerId:      utils.NewRandomInt64(1, 10),
		PlayUrl:      utils.NewRandomString(10),
		CoverUrl:     utils.NewRandomString(10),
		CommentCount: utils.NewRandomInt64(1, 10),
		LikeCount:    utils.NewRandomInt64(1, 10),
		CreateTime:   utils.NewRandomInt64(1, 10),
		UpdateTime:   utils.NewRandomInt64(1, 10),
	}}}
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockSearchRpc.EXPECT().SearchVideo(gomock.Any(), gomock.Any()).Return(resp, nil)

	// 表格驱动测试
	testCases := []struct {
		name string
		req  *types.SearchVideoReq
		err  error
	}{
		{
			name: "search_video_with_validate_error",
			req:  &types.SearchVideoReq{Keyword: "test", Sort: 0},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "search_video_with_es_error",
			req:  &types.SearchVideoReq{Keyword: "test", Sort: 0},
			err:  errors.Wrapf(esError, "req: %v", &types.SearchVideoReq{Keyword: "test", Sort: 0}),
		},
		{
			name: "search_video_success",
			req:  &types.SearchVideoReq{Keyword: "test", Sort: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := videoLogic.Video(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
