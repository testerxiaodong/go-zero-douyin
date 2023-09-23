package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-douyin/apps/user/cmd/rpc/internal/config"
	"go-zero-douyin/apps/user/cmd/rpc/internal/model"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(sqlConn, c.Cache),
	}
}
