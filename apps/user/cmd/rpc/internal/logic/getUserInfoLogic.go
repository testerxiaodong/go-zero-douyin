package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/common/xerr"
	"gorm.io/gorm"

	"go-zero-douyin/apps/user/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/user/cmd/rpc/pb"

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
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "GetUserInfo empty param")
	}
	// 查询用户
	user, err := l.svcCtx.Query.User.WithContext(l.ctx).Where(l.svcCtx.Query.User.ID.Eq(in.GetId())).First()
	// 数据库查询出错处理
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "Get user by id failed, user_id: %d", in.GetId())
	}
	// 用户不存在处理
	if user == nil {
		return nil, errors.Wrapf(ErrUserNotFound, "user_id: %d", in.GetId())
	}
	// 返回数据
	return &pb.GetUserInfoResp{
		User: &pb.UserInfo{
			Id:       user.ID,
			Username: user.Username,
		},
	}, nil
}
