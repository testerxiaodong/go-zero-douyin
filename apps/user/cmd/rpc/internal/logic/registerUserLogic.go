package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/user/cmd/rpc/internal/model"
	"go-zero-douyin/apps/user/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserAlreadyRegister = xerr.NewErrMsg("user has been registered")

type RegisterUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterUserLogic {
	return &RegisterUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterUserLogic) RegisterUser(in *pb.RegisterUserReq) (*pb.RegisterUserResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Register user empty param")
	}
	if len(in.GetUsername()) == 0 || len(in.GetPassword()) == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Register user error param")
	}

	user, err := l.svcCtx.UserModel.FindOneByUsernameIsDelete(l.ctx, in.GetUsername(), xconst.DelStateNo)
	// 查询数据库时出现错误
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "find user by username failed, username: %s, err: %v", in.GetUsername(), err)
	}
	// 用户已存在
	if user != nil {
		return nil, errors.Wrapf(ErrUserAlreadyRegister, "register user exists username: %s", in.GetUsername())
	}
	// 插入用户
	newUser := &model.User{}
	newUser.Username = in.GetUsername()
	// 密码加密
	newUser.Password = utils.Md5ByString(in.GetPassword())
	// 插入用户
	insertResult, err := l.svcCtx.UserModel.Insert(l.ctx, nil, newUser)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "insert user failed, username: %s, password: %s", in.GetUsername(), in.GetPassword())
	}
	uid, _ := insertResult.LastInsertId()
	// 生成token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&pb.GenerateTokenReq{UserId: uid})
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "Generate token failed, user_id: %d", uid)
	}
	return &pb.RegisterUserResp{
		AccessToken:  tokenResp.AccessToken,
		ExpireTime:   tokenResp.ExpireTime,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}
