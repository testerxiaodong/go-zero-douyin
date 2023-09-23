package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SectionModel = (*customSectionModel)(nil)

type (
	// SectionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSectionModel.
	SectionModel interface {
		sectionModel
	}

	customSectionModel struct {
		*defaultSectionModel
	}
)

// NewSectionModel returns a model for the database table.
func NewSectionModel(conn sqlx.SqlConn, c cache.CacheConf) SectionModel {
	return &customSectionModel{
		defaultSectionModel: newSectionModel(conn, c),
	}
}
