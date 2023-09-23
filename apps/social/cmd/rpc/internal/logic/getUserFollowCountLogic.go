package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserFollowCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserFollowCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFollowCountLogic {
	return &GetUserFollowCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserFollowCountLogic) GetUserFollowCount(in *pb.GetUserFollowCountReq) (*pb.GetUserFollowCountResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get user follow count with empty param")
	}
	if in.GetUserId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get user follow count with empty user_id")
	}
	// 查询数据库
	followCount, err := l.svcCtx.FollowCountModel.FindOneByUserId(l.ctx, in.GetUserId())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
			"查询用户关注数粉丝数失败, err: %v, user_id: %d", err, in.GetUserId())
	}
	if followCount == nil {
		return &pb.GetUserFollowCountResp{FollowCount: 0}, nil
	}
	return &pb.GetUserFollowCountResp{FollowCount: followCount.FollowCount}, nil
}
