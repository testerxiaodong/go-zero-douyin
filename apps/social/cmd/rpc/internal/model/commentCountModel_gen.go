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
	commentCountFieldNames          = builder.RawFieldNames(&CommentCount{})
	commentCountRows                = strings.Join(commentCountFieldNames, ",")
	commentCountRowsExpectAutoSet   = strings.Join(stringx.Remove(commentCountFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	commentCountRowsWithPlaceHolder = strings.Join(stringx.Remove(commentCountFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheGoZeroDouyinCommentCountIdPrefix              = "cache:goZeroDouyin:commentCount:id:"
	cacheGoZeroDouyinCommentCountVideoIdIsDeletePrefix = "cache:goZeroDouyin:commentCount:videoId:isDelete:"
)

type (
	commentCountModel interface {
		Insert(ctx context.Context, session sqlx.Session, data *CommentCount) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*CommentCount, error)
		FindOneByVideoIdIsDelete(ctx context.Context, videoId int64, isDelete int64) (*CommentCount, error)
		Update(ctx context.Context, session sqlx.Session, data *CommentCount) (sql.Result, error)
		UpdateWithVersion(ctx context.Context, session sqlx.Session, data *CommentCount) error
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		SelectBuilder() squirrel.SelectBuilder
		DeleteSoft(ctx context.Context, session sqlx.Session, data *CommentCount) error
		FindSum(ctx context.Context, sumBuilder squirrel.SelectBuilder, field string) (float64, error)
		FindCount(ctx context.Context, countBuilder squirrel.SelectBuilder, field string) (int64, error)
		FindAll(ctx context.Context, rowBuilder squirrel.SelectBuilder, orderBy string) ([]*CommentCount, error)
		FindPageListByPage(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*CommentCount, error)
		FindPageListByPageWithTotal(ctx context.Context, rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*CommentCount, int64, error)
		FindPageListByIdDESC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*CommentCount, error)
		FindPageListByIdASC(ctx context.Context, rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*CommentCount, error)
		Delete(ctx context.Context, session sqlx.Session, id int64) error
	}

	defaultCommentCountModel struct {
		sqlc.CachedConn
		table string
	}

	CommentCount struct {
		Id           int64     `db:"id"`            // 评论数id
		VideoId      int64     `db:"video_id"`      // 视频id
		CommentCount int64     `db:"comment_count"` // 评论数
		CreateTime   time.Time `db:"create_time"`   // 创建时间
		UpdateTime   time.Time `db:"update_time"`   // 更新时间
		DeleteTime   time.Time `db:"delete_time"`   // 删除时间
		IsDelete     int64     `db:"is_delete"`     // 是否被删除
		Version      int64     `db:"version"`       // 版本号
	}
)

func newCommentCountModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultCommentCountModel {
	return &defaultCommentCountModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`comment_count`",
	}
}

func (m *defaultCommentCountModel) Insert(ctx context.Context, session sqlx.Session, data *CommentCount) (sql.Result, error) {
	data.DeleteTime = time.Unix(0, 0)
	data.IsDelete = xconst.DelStateNo
	goZeroDouyinCommentCountIdKey := fmt.Sprintf("%s%v", cacheGoZeroDouyinCommentCountIdPrefix, data.Id)
	goZeroDouyinCommentCountVideoIdIsDeleteKey := fmt.Sprintf("%s%v:%v", cacheGoZeroDouyinCommentCountVideoIdIsDeletePrefix, data.VideoId, data.IsDelete)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, commentCountRowsExpectAutoSet)
		if session != nil {
			return session.ExecCtx(ctx, query, data.VideoId, data.CommentCount, data.DeleteTime, data.IsDelete, data.Version)
		}
		return conn.ExecCtx(ctx, query, data.VideoId, data.CommentCount, data.DeleteTime, data.IsDelete, data.Version)
	}, goZeroDouyinCommentCountIdKey, goZeroDouyinCommentCountVideoIdIsDeleteKey)
}

func (m *defaultCommentCountModel) FindOne(ctx context.Context, id int64) (*CommentCount, error) {
	goZeroDouyinCommentCountIdKey := fmt.Sprintf("%s%v", cacheGoZeroDouyinCommentCountIdPrefix, id)
	var resp CommentCount
	err := m.QueryRowCtx(ctx, &resp, goZeroDouyinCommentCountIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and is_delete = ? limit 1", commentCountRows, m.table)
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

func (m *defaultCommentCountModel) FindOneByVideoIdIsDelete(ctx context.Context, videoId int64, isDelete int64) (*CommentCount, error) {
	goZeroDouyinCommentCountVideoIdIsDeleteKey := fmt.Sprintf("%s%v:%v", cacheGoZeroDouyinCommentCountVideoIdIsDeletePrefix, videoId, isDelete)
	var resp CommentCount
	err := m.QueryRowIndexCtx(ctx, &resp, goZeroDouyinCommentCountVideoIdIsDeleteKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `video_id` = ? and `is_delete` = ? and is_delete = ? limit 1", commentCountRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, videoId, isDelete, xconst.DelStateNo); err != nil {
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

func (m *defaultCommentCountModel) Update(ctx context.Context, session sqlx.Session, newData *CommentCount) (sql.Result, error) {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return nil, err
	}
	goZeroDouyinCommentCountIdKey := fmt.Sprintf("%s%v", cacheGoZeroDouyinCommentCountIdPrefix, data.Id)
	goZeroDouyinCommentCountVideoIdIsDeleteKey := fmt.Sprintf("%s%v:%v", cacheGoZeroDouyinCommentCountVideoIdIsDeletePrefix, data.VideoId, data.IsDelete)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, commentCountRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, newData.VideoId, newData.CommentCount, newData.DeleteTime, newData.IsDelete, newData.Version, newData.Id)
		}
		return conn.ExecCtx(ctx, query, newData.VideoId, newData.CommentCount, newData.DeleteTime, newData.IsDelete, newData.Version, newData.Id)
	}, goZeroDouyinCommentCountIdKey, goZeroDouyinCommentCountVideoIdIsDeleteKey)
}

func (m *defaultCommentCountModel) UpdateWithVersion(ctx context.Context, session sqlx.Session, newData *CommentCount) error {

	oldVersion := newData.Version
	newData.Version += 1

	var sqlResult sql.Result
	var err error

	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}
	goZeroDouyinCommentCountIdKey := fmt.Sprintf("%s%v", cacheGoZeroDouyinCommentCountIdPrefix, data.Id)
	goZeroDouyinCommentCountVideoIdIsDeleteKey := fmt.Sprintf("%s%v:%v", cacheGoZeroDouyinCommentCountVideoIdIsDeletePrefix, data.VideoId, data.IsDelete)
	sqlResult, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, commentCountRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, newData.VideoId, newData.CommentCount, newData.DeleteTime, newData.IsDelete, newData.Version, newData.Id, oldVersion)
		}
		return conn.ExecCtx(ctx, query, newData.VideoId, newData.CommentCount, newData.DeleteTime, newData.IsDelete, newData.Version, newData.Id, oldVersion)
	}, goZeroDouyinCommentCountIdKey, goZeroDouyinCommentCountVideoIdIsDeleteKey)
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

func (m *defaultCommentCountModel) DeleteSoft(ctx context.Context, session sqlx.Session, data *CommentCount) error {
	data.IsDelete = data.Id
	data.DeleteTime = time.Now()
	if err := m.UpdateWithVersion(ctx, session, data); err != nil {
		return errors.Wrapf(errors.New("delete soft failed "), "CommentCountModel delete err : %+v", err)
	}
	return nil
}

func (m *defaultCommentCountModel) FindSum(ctx context.Context, builder squirrel.SelectBuilder, field string) (float64, error) {

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

func (m *defaultCommentCountModel) FindCount(ctx context.Context, builder squirrel.SelectBuilder, field string) (int64, error) {

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

func (m *defaultCommentCountModel) FindAll(ctx context.Context, builder squirrel.SelectBuilder, orderBy string) ([]*CommentCount, error) {

	builder = builder.Columns(commentCountRows)

	if orderBy == "" {
		builder = builder.OrderBy("id DESC")
	} else {
		builder = builder.OrderBy(orderBy)
	}

	query, values, err := builder.Where("is_delete = ?", xconst.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*CommentCount
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultCommentCountModel) FindPageListByPage(ctx context.Context, builder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*CommentCount, error) {

	builder = builder.Columns(commentCountRows)

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

	var resp []*CommentCount
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultCommentCountModel) FindPageListByPageWithTotal(ctx context.Context, builder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*CommentCount, int64, error) {

	total, err := m.FindCount(ctx, builder, "id")
	if err != nil {
		return nil, 0, err
	}

	builder = builder.Columns(commentCountRows)

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

	var resp []*CommentCount
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, total, nil
	default:
		return nil, total, err
	}
}

func (m *defaultCommentCountModel) FindPageListByIdDESC(ctx context.Context, builder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*CommentCount, error) {

	builder = builder.Columns(commentCountRows)

	if preMinId > 0 {
		builder = builder.Where(" id < ? ", preMinId)
	}

	query, values, err := builder.Where("is_delete = ?", xconst.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*CommentCount
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultCommentCountModel) FindPageListByIdASC(ctx context.Context, builder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*CommentCount, error) {

	builder = builder.Columns(commentCountRows)

	if preMaxId > 0 {
		builder = builder.Where(" id > ? ", preMaxId)
	}

	query, values, err := builder.Where("is_delete = ?", xconst.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*CommentCount
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultCommentCountModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {

	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})

}

func (m *defaultCommentCountModel) SelectBuilder() squirrel.SelectBuilder {
	return squirrel.Select().From(m.table)
}
func (m *defaultCommentCountModel) Delete(ctx context.Context, session sqlx.Session, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	goZeroDouyinCommentCountIdKey := fmt.Sprintf("%s%v", cacheGoZeroDouyinCommentCountIdPrefix, id)
	goZeroDouyinCommentCountVideoIdIsDeleteKey := fmt.Sprintf("%s%v:%v", cacheGoZeroDouyinCommentCountVideoIdIsDeletePrefix, data.VideoId, data.IsDelete)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.ExecCtx(ctx, query, id)
		}
		return conn.ExecCtx(ctx, query, id)
	}, goZeroDouyinCommentCountIdKey, goZeroDouyinCommentCountVideoIdIsDeleteKey)
	return err
}
func (m *defaultCommentCountModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheGoZeroDouyinCommentCountIdPrefix, primary)
}
func (m *defaultCommentCountModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? and is_delete = ? limit 1", commentCountRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary, xconst.DelStateNo)
}

func (m *defaultCommentCountModel) tableName() string {
	return m.table
}
