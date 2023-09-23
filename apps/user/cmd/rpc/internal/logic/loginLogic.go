package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/user/cmd/rpc/internal/model"
	"go-zero-douyin/apps/user/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUsernamePwdError = xerr.NewErrMsg("账号或密码不正确")

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	// todo: add your logic here and delete this line
	// 校验参数
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "empty param for user login")
	}

	// 查询用户
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.GetUsername())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "根据用户名查询用户失败, err: %v, username: %s", err, in.GetUsername())
	}

	// 用户不存在
	if user == nil {
		return nil, errors.Wrapf(ErrUserNotFound, "username: %s", in.GetUsername())
	}
	// 校验密码
	if user.Password != utils.Md5ByString(in.GetPassword()) {
		return nil, errors.Wrap(ErrUsernamePwdError, "密码不匹配")
	}

	// 生成token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&pb.GenerateTokenReq{
		UserId: user.Id,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "user_id: %d", user.Id)
	}

	// 返回rpc响应
	return &pb.LoginResp{
		AccessToken:  tokenResp.AccessToken,
		RefreshAfter: tokenResp.RefreshAfter,
		ExpireTime:   tokenResp.ExpireTime,
	}, nil
}
