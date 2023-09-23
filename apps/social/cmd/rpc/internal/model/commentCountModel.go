package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CommentCountModel = (*customCommentCountModel)(nil)

type (
	// CommentCountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommentCountModel.
	CommentCountModel interface {
		commentCountModel
	}

	customCommentCountModel struct {
		*defaultCommentCountModel
	}
)

// NewCommentCountModel returns a model for the database table.
func NewCommentCountModel(conn sqlx.SqlConn, c cache.CacheConf) CommentCountModel {
	return &customCommentCountModel{
		defaultCommentCountModel: newCommentCountModel(conn, c),
	}
}
