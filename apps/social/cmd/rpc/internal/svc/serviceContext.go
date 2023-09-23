package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-douyin/apps/social/cmd/rpc/internal/config"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
)

type ServiceContext struct {
	Config            config.Config
	CommentModel      model.CommentModel
	CommentCountModel model.CommentCountModel
	FollowModel       model.FollowModel
	FollowCountModel  model.FollowCountModel
	LikeModel         model.LikeModel
	LikeCountModel    model.LikeCountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:            c,
		CommentModel:      model.NewCommentModel(sqlConn, c.Cache),
		CommentCountModel: model.NewCommentCountModel(sqlConn, c.Cache),
		FollowModel:       model.NewFollowModel(sqlConn, c.Cache),
		FollowCountModel:  model.NewFollowCountModel(sqlConn, c.Cache),
		LikeModel:         model.NewLikeModel(sqlConn, c.Cache),
		LikeCountModel:    model.NewLikeCountModel(sqlConn, c.Cache),
	}
}
