package user

import (
	"context"
	"github.com/pkg/errors"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/user/cmd/api/internal/svc"
	"go-zero-douyin/apps/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}

	// 调用userRpc获取用户信息
	userInfoResp, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &pb.GetUserInfoReq{Id: req.Id})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}
	// 调用socialRpc获取用户粉丝数
	followerCountResp, err := l.svcCtx.SocialRpc.GetUserFollowerCount(l.ctx, &socialPb.GetUserFollowerCountReq{UserId: req.Id})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}
	// 调用socialRpc获取用户关注数
	followCountResp, err := l.svcCtx.SocialRpc.GetUserFollowCount(l.ctx, &socialPb.GetUserFollowCountReq{UserId: req.Id})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}

	return &types.UserInfoResp{
		User: types.UserInfo{
			Id:            userInfoResp.User.Id,
			Username:      userInfoResp.User.Username,
			FollowerCount: followerCountResp.FollowerCount,
			FollowCount:   followCountResp.FollowCount,
		},
	}, nil
}
