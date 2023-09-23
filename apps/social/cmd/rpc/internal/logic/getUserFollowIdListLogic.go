package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserFollowIdListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserFollowIdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFollowIdListLogic {
	return &GetUserFollowIdListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserFollowIdListLogic) GetUserFollowIdList(in *pb.GetUserFollowIdListReq) (*pb.GetUserFollowIdListResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follow id list with empty param")
	}
	if in.GetUserId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follow id list with empty user_id")
	}
	if in.GetPage() > 5 {
		return nil, errors.Wrap(xerr.NewErrMsg("系统不允许超过五页"), "关注列表业务校验")
	}
	builder := l.svcCtx.FollowModel.SelectBuilder().Where(squirrel.Eq{"follower_id": in.GetUserId()}).Where(squirrel.Eq{"status": xconst.FollowStateYes})
	follows, err := l.svcCtx.FollowModel.FindPageListByPage(l.ctx, builder, in.GetPage(), in.GetPageSize(), "create_time DESC")
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
			"查询用户关注列表失败, err: %v, user_id: %d", err, in.GetUserId())
	}
	if len(follows) == 0 {
		return &pb.GetUserFollowIdListResp{}, nil
	}
	ids := make([]int64, 0)
	// 拼接响应
	for _, follow := range follows {
		ids = append(ids, follow.UserId)
	}
	return &pb.GetUserFollowIdListResp{UserIdList: ids}, nil
}
