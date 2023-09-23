package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnfollowUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnfollowUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfollowUserLogic {
	return &UnfollowUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnfollowUserLogic) UnfollowUser(in *pb.UnfollowUserReq) (*pb.UnfollowUserResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unfollow user with empty param")
	}
	if in.GetUserId() == 0 || in.GetFollowerId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "unfollow user with empty follower_id or user_id")
	}
	// 不能对自己操作
	if in.GetUserId() == in.GetFollowerId() {
		return nil, errors.Wrapf(xerr.NewErrMsg("不能对自己操作"), "req: %v", in)
	}
	// 查询关注记录
	record, err := l.svcCtx.FollowModel.FindOneByUserIdFollowerId(l.ctx, in.GetUserId(), in.GetFollowerId())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
			"查询关注记录失败, err: %v, user_id: %v, follower_id: %d", err, in.GetUserId(), in.GetFollowerId())
	}
	// 没有记录，返回错误信息
	if record == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("未关注过该用户"),
			"user_id: %v, follower_id: %d", in.GetUserId(), in.GetFollowerId())
	}
	// 记录存在且为未关注状态，直接返回
	if record != nil && record.Status == xconst.FollowStateNo {
		return &pb.UnfollowUserResp{}, nil
	}
	// 记录存在且为已关注状态，修改状态，更新粉丝数/关注数
	if record != nil && record.Status == xconst.FollowStateYes {
		if err := l.svcCtx.FollowModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {
			// 修改状态
			record.Status = xconst.FollowStateNo
			err := l.svcCtx.FollowModel.UpdateWithVersion(l.ctx, session, record)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR),
					"更新关注状态失败, err: %v, user_id: %d, follower_id: %d", err, in.GetUserId(), in.GetFollowerId())
			}
			// 更新粉丝数：-1
			followerCount, err := l.svcCtx.FollowCountModel.FindOneByUserId(l.ctx, in.GetUserId())
			if err != nil && !errors.Is(err, model.ErrNotFound) {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
					"查询用户粉丝/关注数失败, err: %v, user_id: %d", err, in.GetUserId())
			}
			followerCount.FollowerCount -= 1
			err = l.svcCtx.FollowCountModel.UpdateWithVersion(l.ctx, session, followerCount)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR),
					"更新用户粉丝数失败, err: %v, user_id: %d", err, in.GetUserId())
			}
			// 更新关注数：-1
			followCount, err := l.svcCtx.FollowCountModel.FindOneByUserId(l.ctx, in.GetFollowerId())
			if err != nil && !errors.Is(err, model.ErrNotFound) {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
					"查询用户粉丝/关注数失败, err: %v, user_id: %d", err, in.GetFollowerId())
			}
			followCount.FollowCount -= 1
			err = l.svcCtx.FollowCountModel.UpdateWithVersion(l.ctx, session, followCount)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR),
					"更新用户关注数失败, err: %v, user_id: %d", err, in.GetFollowerId())
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return &pb.UnfollowUserResp{}, nil
}
