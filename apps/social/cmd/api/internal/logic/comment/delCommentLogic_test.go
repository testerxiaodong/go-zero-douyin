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
	gloablMock "go-zero-douyin/mock"
	"testing"
)

func TestDelCommentLogic_DelComment(t *testing.T) {
	ctl := gomock.NewController(t)

	mockValidator := gloablMock.NewMockValidator(ctl)

	mockSocialRpc := socialMock.NewMockSocial(ctl)

	serviceContext := &svc.ServiceContext{Validator: mockValidator, SocialRpc: mockSocialRpc}

	delCommentLogic := comment.NewDelCommentLogic(context.Background(), serviceContext)

	// 参数校验失败mock
	validateResult := utils.NewRandomString(10)
	mockValidator.EXPECT().Validate(gomock.Any()).Return(validateResult)

	// 参数校验成功，但SocialRpc.DelComment失败mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	socialRpcError := errors.New("SocialRpc.DelComment error")
	mockSocialRpc.EXPECT().DelComment(gomock.Any(), gomock.Any()).Return(&socialPb.DelCommentResp{}, socialRpcError)

	// 参数校验成功，且SocialRpc.DelComment调用成功的mock
	mockValidator.EXPECT().Validate(gomock.Any()).Return("")
	mockSocialRpc.EXPECT().DelComment(gomock.Any(), gomock.Any()).Return(&socialPb.DelCommentResp{}, nil)

	testCases := []struct {
		name string
		req  *types.DelCommentReq
		err  error
	}{
		{
			name: "del_comment_with_validate_error",
			req:  &types.DelCommentReq{CommentId: 1},
			err:  xerr.NewErrMsg(validateResult),
		},
		{
			name: "del_comment_with_social_rpc_error",
			req:  &types.DelCommentReq{CommentId: 1},
			err:  errors.Wrapf(socialRpcError, "req: %v", &types.DelCommentReq{CommentId: 1}),
		},
		{
			name: "del_comment_success",
			req:  &types.DelCommentReq{CommentId: 1},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := delCommentLogic.DelComment(tc.req)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
