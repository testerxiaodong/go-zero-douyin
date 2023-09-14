package user

import (
	"context"
	"github.com/pkg/errors"
	searchPb "go-zero-douyin/apps/search/cmd/rpc/pb"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/user/cmd/api/internal/svc"
	"go-zero-douyin/apps/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncUserToEsByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncUserToEsByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncUserToEsByIdLogic {
	return &SyncUserToEsByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncUserToEsByIdLogic) SyncUserToEsById(req *types.SyncUserToEsByIdReq) error {
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}
	// 调用userRpc
	userInfoResp, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &pb.GetUserInfoReq{Id: req.UserId})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	// 调用socialRpc获取用户粉丝数
	followerCountResp, err := l.svcCtx.SocialRpc.GetUserFollowerCount(l.ctx, &socialPb.GetUserFollowerCountReq{UserId: req.UserId})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	// 调用socialRpc获取用户关注数
	followCountResp, err := l.svcCtx.SocialRpc.GetUserFollowCount(l.ctx, &socialPb.GetUserFollowCountReq{UserId: req.UserId})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	// 调用searchRpc同步用户信息
	_, err = l.svcCtx.SearchRpc.SyncUserInfo(l.ctx, &searchPb.SyncUserInfoReq{
		User: &searchPb.User{
			Id:            userInfoResp.User.Id,
			Username:      userInfoResp.User.Username,
			FollowerCount: followerCountResp.FollowerCount,
			FollowCount:   followCountResp.FollowCount,
		},
	})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	return nil
}
