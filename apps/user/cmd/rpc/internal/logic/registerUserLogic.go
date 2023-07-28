package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/user/cmd/rpc/internal/model"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"gorm.io/gorm"

	"go-zero-douyin/apps/user/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserAlreadyRegisterError = xerr.NewErrMsg("user has been registered")

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
	user, err := l.svcCtx.Query.User.WithContext(l.ctx).Where(l.svcCtx.Query.User.Username.Eq(in.GetUsername())).First()
	// 查询数据库时出现错误
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "find user by username failed, username: %s, err: %v", in.GetUsername(), err)
	}
	// 用户已存在
	if user != nil {
		return nil, errors.Wrapf(ErrUserAlreadyRegisterError, "register user exists username: %s", in.GetUsername())
	}
	// 插入用户
	u := &model.User{}
	// 深拷贝
	err = copier.Copy(u, in)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "复制db -> pb失败")
	}
	// 密码加密
	u.Password = utils.Md5ByString(in.GetPassword())
	err = l.svcCtx.Query.User.WithContext(l.ctx).Create(u)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "insert user failed, username: %s, password: %s", in.GetUsername(), in.GetPassword())
	}
	// 生成token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&pb.GenerateTokenReq{UserId: u.ID})
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "Generate token failed, user_id: %d", u.ID)
	}
	return &pb.RegisterUserResp{
		AccessToken:  tokenResp.AccessToken,
		ExpireTime:   tokenResp.ExpireTime,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}
