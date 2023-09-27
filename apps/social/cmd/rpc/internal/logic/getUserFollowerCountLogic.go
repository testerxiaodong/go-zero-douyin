package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserFollowerCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserFollowerCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFollowerCountLogic {
	return &GetUserFollowerCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserFollowerCountLogic) GetUserFollowerCount(in *pb.GetUserFollowerCountReq) (*pb.GetUserFollowerCountResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follower count with empty param")
	}
	if in.GetUserId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follower count with empty user_id")
	}
	// 查询数据库
	followCount, err := l.svcCtx.FollowCountModel.FindOneByUserIdIsDelete(l.ctx, in.GetUserId(), xconst.DelStateNo)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
			"查询用户粉丝数关注数失败, err: %v, user_id: %d", err, in.GetUserId())
	}
	if followCount == nil {
		return &pb.GetUserFollowerCountResp{FollowerCount: 0}, nil
	}
	return &pb.GetUserFollowerCountResp{FollowerCount: followCount.FollowerCount}, nil
}
