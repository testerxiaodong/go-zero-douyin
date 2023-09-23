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

type GetUserLikeVideoIdListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLikeVideoIdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLikeVideoIdListLogic {
	return &GetUserLikeVideoIdListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLikeVideoIdListLogic) GetUserLikeVideoIdList(in *pb.GetUserLikeVideoIdListReq) (*pb.GetUserLikeVideoIdListResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user like video id list with empty param")
	}
	if in.GetUserId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user like video id list with empty user_id")
	}
	builder := l.svcCtx.LikeModel.SelectBuilder().Where(squirrel.Eq{"user_id": in.GetUserId()}).Where(squirrel.Eq{"status": xconst.LikeStateYes})
	likes, err := l.svcCtx.LikeModel.FindPageListByPage(l.ctx, builder, in.GetPage(), in.GetPageSize(), "create_time DESC")
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
			"查询用户点赞视频列表失败, err: %v, user_id: %d", err, in.GetUserId())
	}
	if len(likes) == 0 {
		return &pb.GetUserLikeVideoIdListResp{}, nil
	}
	ids := make([]int64, 0)
	// 拼接响应
	for _, like := range likes {
		ids = append(ids, like.VideoId)
	}
	return &pb.GetUserLikeVideoIdListResp{VideoIdList: ids}, nil
}
