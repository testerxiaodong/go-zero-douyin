package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/user/cmd/rpc/internal/model"
	"go-zero-douyin/common/xerr"
	"gorm.io/gorm"

	"go-zero-douyin/apps/user/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/user/cmd/rpc/pb"

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

	// 查询用户信息
	userRecord, err := l.svcCtx.UserDo.GetUserById(l.ctx, in.GetId())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "Find user by id failed, user_id: %d", in.GetId())
	}
	if userRecord == nil {
		return nil, errors.Wrapf(ErrUserNotFound, "id: %d", in.GetId())
	}

	// 更新数据
	user := &model.User{}
	err = copier.Copy(user, in)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "copy user update info failed, info: %v", in)
	}
	_, err = l.svcCtx.UserDo.UpdateUserInfo(l.ctx, user, in.GetId())
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR), "update mysql user info failed, user info: %v", user)
	}
	return &pb.UpdateUserResp{}, nil
}
