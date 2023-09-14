package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/user/cmd/rpc/internal/model"
	"go-zero-douyin/common/message"
	"go-zero-douyin/common/utils"
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

	user := &model.User{}
	// 查询用户信息
	if len(in.GetUsername()) > 0 {
		userRecord, err := l.svcCtx.UserDo.GetUserByUsername(l.ctx, in.GetUsername())
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "Find user by id failed, user_id: %d, err: %v", in.GetId(), err)
		}
		if userRecord != nil && userRecord.ID != in.GetId() {
			return nil, errors.Wrapf(ErrUserAlreadyRegister, "username: %s", in.GetUsername())
		}
		user.Username = in.GetUsername()
	}
	if len(in.GetPassword()) > 0 {
		user.Password = utils.Md5ByString(in.GetPassword())
	}
	// 更新数据
	_, err := l.svcCtx.UserDo.UpdateUserInfo(l.ctx, user, in.GetId())
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_UPDATE_ERR), "update mysql user info failed, err: %v", err)
	}
	// 发布更新es用户文档的消息
	msg, _ := json.Marshal(message.MysqlUserUpdateMessage{UserId: in.GetId()})
	err = l.svcCtx.Rabbit.Send("", "MysqlUserUpdateMq", msg)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_UPDATE_ERR), "发布更新es用户文档信息失败, err: %v", err)
	}
	return &pb.UpdateUserResp{}, nil
}
