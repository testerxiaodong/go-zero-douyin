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

type FollowUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowUserLogic {
	return &FollowUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FollowUser 关注功能
func (l *FollowUserLogic) FollowUser(in *pb.FollowUserReq) (*pb.FollowUserResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "follow user with empty param")
	}
	if in.GetUserId() == 0 || in.GetFollowerId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR),
			"follow user with empty follower_id or user_id")
	}

	// 用户不能关注自己
	if in.GetUserId() == in.GetFollowerId() {
		return nil, errors.Wrapf(xerr.NewErrMsg("不能对自己操作"), "req: %v", in)
	}

	// 查询数据库
	follow, err := l.svcCtx.FollowModel.FindOneByUserIdFollowerIdIsDelete(l.ctx, in.GetUserId(), in.GetFollowerId(), xconst.DelStateNo)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
			"search user is alreaddy follow user from db failed, err: %v, follower_id: %d user_id: %d",
			err, in.GetFollowerId(), in.GetUserId())
	}
	// 有记录，且关注状态为已关注，直接返回
	if follow != nil && follow.Status == xconst.FollowStateYes {
		return &pb.FollowUserResp{}, nil
	}
	// 有记录，且关注状态为已取消关注，关注状态更新，关注数粉丝数更新
	if follow != nil && follow.Status == xconst.FollowStateNo {
		if err := l.svcCtx.FollowModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {
			// 更新关注状态
			follow.Status = xconst.FollowStateYes
			err := l.svcCtx.FollowModel.UpdateWithVersion(l.ctx, session, follow)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR),
					"更新关注状态失败, err: %v, user_id: %d, follower_id: %d", err, in.GetUserId(), in.GetFollowerId())
			}
			// 查询粉丝数与关注数
			followerCountRecord, err := l.svcCtx.FollowCountModel.FindOneByUserIdIsDelete(l.ctx, in.GetUserId(), xconst.DelStateNo)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
					"查询用户粉丝数失败, err: %v, user_id: %d", err, in.GetUserId())
			}
			followerCountRecord.FollowerCount += 1
			err = l.svcCtx.FollowCountModel.UpdateWithVersion(l.ctx, session, followerCountRecord)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR),
					"更新用户粉丝数失败, err: %v, user_id: %v", err, in.GetUserId())
			}
			followCountRecord, err := l.svcCtx.FollowCountModel.FindOneByUserIdIsDelete(l.ctx, in.GetFollowerId(), xconst.DelStateNo)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
					"查询用户关注数失败, err: %v, user_id: %v", err, in.GetFollowerId())
			}
			followCountRecord.FollowCount += 1
			err = l.svcCtx.FollowCountModel.UpdateWithVersion(l.ctx, session, followCountRecord)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR),
					"更新用户关注数失败, err: %v, user_id: %v", err, in.GetFollowerId())
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	// 没有记录，插入记录，更新关注数粉丝数
	if follow == nil {
		if err := l.svcCtx.FollowModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {
			// 插入记录
			newFollow := &model.Follow{}
			newFollow.FollowerId = in.GetFollowerId()
			newFollow.UserId = in.GetUserId()
			newFollow.Status = xconst.FollowStateYes
			_, err := l.svcCtx.FollowModel.Insert(l.ctx, session, newFollow)
			if err != nil {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR),
					"插入关注记录失败, err: %v, user_id: %d, follower_id: %d", err, in.GetUserId(), in.GetFollowerId())
			}
			// 更新被关注用户粉丝数
			followerCountRecord, err := l.svcCtx.FollowCountModel.FindOneByUserIdIsDelete(l.ctx, in.GetUserId(), xconst.DelStateNo)
			if err != nil && !errors.Is(err, model.ErrNotFound) {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
					"根据用户id查询粉丝关注数失败, err: %v, user_id: %v", err, in.GetUserId())
			}
			if followerCountRecord == nil {
				// 粉丝数关注数记录不存在，则插入一条
				followCount := &model.FollowCount{}
				followCount.UserId = in.GetUserId()
				followCount.FollowCount = 0
				followCount.FollowerCount = 1
				_, err = l.svcCtx.FollowCountModel.Insert(l.ctx, session, followCount)
				if err != nil {
					return errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "插入用户粉丝数关注数记录失败, err: %v, user_id: %v", err, in.GetUserId())
				}
			} else {
				// 粉丝数关注数记录存在，更新
				followerCountRecord.FollowerCount += 1
				err = l.svcCtx.FollowCountModel.UpdateWithVersion(l.ctx, session, followerCountRecord)
				if err != nil {
					return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR), "更新用户粉丝数失败, err: %v, user_id: %d", err, in.GetUserId())
				}
			}
			// 更新关注用户关注数
			followCountRecord, err := l.svcCtx.FollowCountModel.FindOneByUserIdIsDelete(l.ctx, in.GetFollowerId(), xconst.DelStateNo)
			if err != nil && !errors.Is(err, model.ErrNotFound) {
				return errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
					"根据用户id查询粉丝关注数失败, err: %v, user_id: %v", err, in.GetFollowerId())
			}
			if followCountRecord == nil {
				// 粉丝数关注数记录不存在，则插入一条
				followCount := &model.FollowCount{}
				followCount.UserId = in.GetFollowerId()
				followCount.FollowCount = 1
				followCount.FollowerCount = 0
				_, err = l.svcCtx.FollowCountModel.Insert(l.ctx, session, followCount)
				if err != nil {
					return errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "插入用户粉丝数关注数记录失败, err: %v, user_id: %v", err, in.GetFollowerId())
				}
			} else {
				// 粉丝数关注数记录存在，更新
				followCountRecord.FollowCount += 1
				err = l.svcCtx.FollowCountModel.UpdateWithVersion(l.ctx, session, followCountRecord)
				if err != nil {
					return errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR), "更新用户粉丝数失败, err: %v, user_id: %d", err, in.GetFollowerId())
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return &pb.FollowUserResp{}, nil
}
