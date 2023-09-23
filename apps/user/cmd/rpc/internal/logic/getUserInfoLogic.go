package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/user/cmd/rpc/internal/model"
	"go-zero-douyin/apps/user/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserNotFound = xerr.NewErrMsg("用户不存在")

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "获取用户信息时参数为nil")
	}
	if in.GetId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "获取用户信息时id为空")
	}
	// 查询用户信息
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.GetId())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "从数据库查询用户信息失败 user_id:%d , err:%v", in.GetId(), err)
	}
	if user == nil {
		return nil, errors.Wrapf(ErrUserNotFound, "id:%d", in.GetId())
	}
	// 返回数据
	return &pb.GetUserInfoResp{
		User: &pb.UserInfo{
			Id:       user.Id,
			Username: user.Username,
		},
	}, nil
}
