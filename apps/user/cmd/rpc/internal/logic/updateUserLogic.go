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

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *pb.UpdateUserReq) (*pb.UpdateUserResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Update user empty param")
	}
	if in.GetId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Updtae user empty user id")
	}
	oldUser, err := l.svcCtx.UserModel.FindOne(l.ctx, in.GetId())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
			"根据id查询用户失败, err: %v, user_id: %d", err, in.GetId())
	}
	if oldUser == nil {
		return nil, errors.Wrapf(ErrUserNotFound, "user_id: %d", in.GetId())
	}
	// 查询用户
	if len(in.GetUsername()) > 0 {
		user, err := l.svcCtx.UserModel.FindOneByUsernameIsDelete(l.ctx, in.GetUsername(), xconst.DelStateNo)
		if err != nil && !errors.Is(err, model.ErrNotFound) {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
				"根据用户名查询用户失败, err: %v, username: %s", err, in.GetUsername())
		}
		if user != nil && user.Id != in.GetId() {
			return nil, errors.Wrapf(ErrUserAlreadyRegister, "username: %s", in.GetUsername())
		}
		oldUser.Username = in.GetUsername()
	}
	if len(in.GetPassword()) > 0 {
		oldUser.Password = utils.Md5ByString(in.GetPassword())
	}
	// 更新数据
	err = l.svcCtx.UserModel.UpdateWithVersion(l.ctx, nil, oldUser)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR), "update mysql user info failed, err: %v", err)
	}
	return &pb.UpdateUserResp{}, nil
}
