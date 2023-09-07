package comment_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go-zero-douyin/apps/social/cmd/api/internal/logic/comment"
	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"
	socialMock "go-zero-douyin/apps/social/cmd/rpc/mock"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"go-zero-douyin/mock"
	"testing"
)

func TestVideoCommentListLogic_VideoCommentList(t *testing.T) {
	ctl := gomock.NewController(t)

	mockValidator := mock.NewMockValidator(ctl)

	mockSocialRpc := socialMock.NewMockSocial(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, SocialRpc: mockSocialRpc}

	videoCommentListLogic := comment.NewVideoCommentListLogic(context.Background(), serviceContext)

	// 参数校验失败mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// 参数校验成功，但SocialRpc.GetVideoCommentListById失败mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	socialRpcError := errors.New("SocialRpc.GetVideoCommentListById error")
	mockSocialRpc.EXPECT().GetVideoCommentListById(gomock.Any(), gomock.Any()).Return(&socialPb.GetCommentListByIdResp{}, socialRpcError)

	// 参数校验成功，且SocialRpc.GetVideoCommentListById调用成功的mock
	expectedResp := &socialPb.GetCommentListByIdResp{
		Comments: []*socialPb.Comment{&socialPb.Comment{Id: 1, UserId: 1, VideoId: 1, Content: "test"},
			&socialPb.Comment{Id: 2, UserId: 2, VideoId: 1, Content: "test"}}}
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockSocialRpc.EXPECT().GetVideoCommentListById(gomock.Any(), gomock.Any()).Return(expectedResp, nil)

	testCases := []struct {
		name string
		req  *types.GetVideoCommentListReq
		err  error
	}{
		{
			name: "video_comment_list_with_validate_error",
			req:  &types.GetVideoCommentListReq{VideoId: 1},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "video_comment_list_with_social_rpc_error",
			req:  &types.GetVideoCommentListReq{VideoId: 1},
			err:  errors.Wrapf(socialRpcError, "req: %v", &types.GetVideoCommentListReq{VideoId: 1}),
		},
		{
			name: "video_comment_list_success",
			req:  &types.GetVideoCommentListReq{VideoId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := videoCommentListLogic.VideoCommentList(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
				assert.Equal(t, len(resp.Comments), 2)
			}
		})
	}
}
