// Code generated by goctl. DO NOT EDIT!

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
	"go-zero-douyin/common/xconst"
)

var (
	sectionFieldNames          = builder.RawFieldNames(&Section{})
	sectionRows                = strings.Join(sectionFieldNames, ",")
	sectionRowsExpectAutoSet   = strings.Join(stringx.Remove(sectionFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	sectionRowsWithPlaceHolder = strings.Join(stringx.Remove(sectionFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheGoZeroDouyinSectionIdPrefix           = "cache:goZeroDouyin:section:id:"
	cacheGoZeroDouyinSectionNameIsDeletePrefix = "cache:goZeroDouyin:section:name:isDelete:"
)

type (
	sectionModel interface {
		Insert(ctx context.Context, session sqlx.Session, data *Section) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Section, error)
		FindOneByNameIsDelete(ctx context.Context, name string, isDelete int64) (*Section, error)
		Update(ctx context.Context, session sqlx.Session, data *Section) (sql.Result, error)
		UpdateWithVersion(ctx context.Context, session sqlx.Session, data *Section) error
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		SelectBuilder() squirrel.SelectBuilder
		DeleteSoft(ctx context.Context, session sqlx.Session, data *Section) error
		FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder, field string) (float64, error)
		FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder, field string) (int64, error)
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder, orderBy string) ([]*Section, error)
		FindPageListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Section, error)
		FindPageListByPageWithTotal(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Section, int64, error)
		FindPageListByIdDESC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Section, error)
		FindPageListByIdASC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*Section, error)
		Delete(ctx context.Context, session sqlx.Session, id int64) error
	}

	defaultSectionModel struct {
		sqlc.CachedConn
		table string
	}

	Section struct {
		Id         int64     `db:"id"`          // 分区id
		Name       string    `db:"name"`        // 分区名
		CreateTime time.Time `db:"create_time"` // 创建时间
		UpdateTime time.Time `db:"update_time"` // 更新时间
		DeleteTime time.Time `db:"delete_time"` // 删除时间
		IsDelete   int64     `db:"is_delete"`   // 是否被删除
		Version    int64     `db:"version"`     // 版本号
	}
)

func newSectionModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultSectionModel {
	return &defaultSectionModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`section`",
	}
}

func (m *defaultSectionModel) Insert(ctx context.Context, session sqlx.Session, data *Section) (sql.Result, error) {
	data.DeleteTime = time.Unix(0, 0)
	data.IsDelete = xconst.DelStateNo
	goZeroDouyinSectionIdKey := fmt.Sprintf("%s%v", cacheGoZeroDouyinSectionIdPrefix, data.Id)
	goZeroDouyinSectionNameIsDeleteKey := fmt.Sprintf("%s%v:%v", cacheGoZeroDouyinSectionNameIsDeletePrefix, data.Name, data.IsDelete)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, sectionRowsExpectAutoSet)
		if session != nil {
			return session.ExecCtx(ctx, query, data.Name, data.DeleteTime, data.IsDelete, data.Version)
		}
		return conn.ExecCtx(ctx, query, data.Name, data.DeleteTime, data.IsDelete, data.Version)
	}, goZeroDouyinSectionIdKey, goZeroDouyinSectionNameIsDeleteKey)
}

func (m *defaultSectionModel) FindOne(ctx context.Context, id int64) (*Section, error) {
	goZeroDouyinSectionIdKey := fmt.Sprintf("%s%v", cacheGoZeroDouyinSectionIdPrefix, id)
	var resp Section
	err := m.QueryRowCtx(ctx, &resp, goZeroDouyinSectionIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and is_delete = ? limit 1", sectionRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id, xconst.DelStateNo)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultSectionModel) FindOneByNameIsDelete(ctx context.Context, name string, isDelete int64) (*Section, error) {
	goZeroDouyinSectionNameIsDeleteKey := fmt.Sprintf("%s%v:%v", cacheGoZeroDouyinSectionNameIsDeletePrefix, name, isDelete)
	var resp Section
	err := m.QueryRowIndexCtx(ctx, &resp, goZeroDouyinSectionNameIsDeleteKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `name` = ? and `is_delete` = ? and is_delete = ? limit 1", sectionRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, name, isDelete, xconst.DelStateNo); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultSectionModel) Update(ctx context.Context, session sqlx.Session, newData *Section) (sql.Result, error) {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return nil, err
	}
	goZeroDouyinSectionIdKey := fmt.Sprintf("%s%v", cacheGoZeroDouyinSectionIdPrefix, data.Id)
	goZeroDouyinSectionNameIsDeleteKey := fmt.Sprintf("%s%v:%v", cacheGoZeroDouyinSectionNameIsDeletePrefix, data.Name, data.IsDelete)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, sectionRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, newData.Name, newData.DeleteTime, newData.IsDelete, newData.Version, newData.Id)
		}
		return conn.ExecCtx(ctx, query, newData.Name, newData.DeleteTime, newData.IsDelete, newData.Version, newData.Id)
	}, goZeroDouyinSectionIdKey, goZeroDouyinSectionNameIsDeleteKey)
}

func (m *defaultSectionModel) UpdateWithVersion(ctx context.Context, session sqlx.Session, newData *Section) error {

	oldVersion := newData.Version
	newData.Version += 1

	var sqlResult sql.Result
	var err error

	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}
	goZeroDouyinSectionIdKey := fmt.Sprintf("%s%v", cacheGoZeroDouyinSectionIdPrefix, data.Id)
	goZeroDouyinSectionNameIsDeleteKey := fmt.Sprintf("%s%v:%v", cacheGoZeroDouyinSectionNameIsDeletePrefix, data.Name, data.IsDelete)
	sqlResult, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, sectionRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, newData.Name, newData.DeleteTime, newData.IsDelete, newData.Version, newData.Id, oldVersion)
		}
		return conn.ExecCtx(ctx, query, newData.Name, newData.DeleteTime, newData.IsDelete, newData.Version, newData.Id, oldVersion)
	}, goZeroDouyinSectionIdKey, goZeroDouyinSectionNameIsDeleteKey)
	if err != nil {
		return err
	}
	updateCount, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}
	if updateCount == 0 {
		return ErrNoRowsUpdate
	}

	return nil
}

func (m *defaultSectionModel) DeleteSoft(ctx context.Context, session sqlx.Session, data *Section) error {
	data.IsDelete = data.Id
	data.DeleteTime = time.Now()
	if err := m.UpdateWithVersion(ctx, session, data); err != nil {
		return errors.Wrapf(errors.New("delete soft failed "), "SectionModel delete err : %+v", err)
	}
	return nil
}

func (m *defaultSectionModel) FindSum(ctx context.Context, builder squirrel.SelectBuilder, field string) (float64, error) {

	if len(field) == 0 {
		return 0, errors.Wrapf(errors.New("FindSum Least One Field"), "FindSum Least One Field")
	}

	builder = builder.Columns("IFNULL(SUM(" + field + "),0)")

	query, values, err := builder.Where("is_delete = ?", xconst.DelStateNo).ToSql()
	if err != nil {
		return 0, err
	}

	var resp float64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultSectionModel) FindCount(ctx context.Context, builder squirrel.SelectBuilder, field string) (int64, error) {

	if len(field) == 0 {
		return 0, errors.Wrapf(errors.New("FindCount Least One Field"), "FindCount Least One Field")
	}

	builder = builder.Columns("COUNT(" + field + ")")

	query, values, err := builder.Where("is_delete = ?", xconst.DelStateNo).ToSql()
	if err != nil {
		return 0, err
	}

	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultSectionModel) FindAll(ctx context.Context, builder squirrel.SelectBuilder, orderBy string) ([]*Section, error) {

	builder = builder.Columns(sectionRows)

	if orderBy == "" {
		builder = builder.OrderBy("id DESC")
	} else {
		builder = builder.OrderBy(orderBy)
	}

	query, values, err := builder.Where("is_delete = ?", xconst.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Section
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultSectionModel) FindPageListByPage(ctx context.Context, builder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Section, error) {

	builder = builder.Columns(sectionRows)

	if orderBy == "" {
		builder = builder.OrderBy("id DESC")
	} else {
		builder = builder.OrderBy(orderBy)
	}

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query, values, err := builder.Where("is_delete = ?", xconst.DelStateNo).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Section
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultSectionModel) FindPageListByPageWithTotal(ctx context.Context, builder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Section, int64, error) {

	total, err := m.FindCount(ctx, builder, "id")
	if err != nil {
		return nil, 0, err
	}

	builder = builder.Columns(sectionRows)

	if orderBy == "" {
		builder = builder.OrderBy("id DESC")
	} else {
		builder = builder.OrderBy(orderBy)
	}

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query, values, err := builder.Where("is_delete = ?", xconst.DelStateNo).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, total, err
	}

	var resp []*Section
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, total, nil
	default:
		return nil, total, err
	}
}

func (m *defaultSectionModel) FindPageListByIdDESC(ctx context.Context, builder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Section, error) {

	builder = builder.Columns(sectionRows)

	if preMinId > 0 {
		builder = builder.Where(" id < ? ", preMinId)
	}

	query, values, err := builder.Where("is_delete = ?", xconst.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Section
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultSectionModel) FindPageListByIdASC(ctx context.Context, builder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*Section, error) {

	builder = builder.Columns(sectionRows)

	if preMaxId > 0 {
		builder = builder.Where(" id > ? ", preMaxId)
	}

	query, values, err := builder.Where("is_delete = ?", xconst.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Section
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultSectionModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {

	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})

}

func (m *defaultSectionModel) SelectBuilder() squirrel.SelectBuilder {
	return squirrel.Select().From(m.table)
}
func (m *defaultSectionModel) Delete(ctx context.Context, session sqlx.Session, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	goZeroDouyinSectionIdKey := fmt.Sprintf("%s%v", cacheGoZeroDouyinSectionIdPrefix, id)
	goZeroDouyinSectionNameIsDeleteKey := fmt.Sprintf("%s%v:%v", cacheGoZeroDouyinSectionNameIsDeletePrefix, data.Name, data.IsDelete)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.ExecCtx(ctx, query, id)
		}
		return conn.ExecCtx(ctx, query, id)
	}, goZeroDouyinSectionIdKey, goZeroDouyinSectionNameIsDeleteKey)
	return err
}
func (m *defaultSectionModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheGoZeroDouyinSectionIdPrefix, primary)
}
func (m *defaultSectionModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? and is_delete = ? limit 1", sectionRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary, xconst.DelStateNo)
}

func (m *defaultSectionModel) tableName() string {
	return m.table
}
