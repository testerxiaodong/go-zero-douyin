package dao

import (
	"context"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"time"
)

type CommentDo interface {
	GetCommentById(ctx context.Context, userId int64) (*model.Comment, error)
	GetCommentListByVideoId(ctx context.Context, videoId int64) ([]*model.Comment, error)
	GetCommentCountByVideoId(ctx context.Context, videoId int64) (int64, error)
	InsertComment(ctx context.Context, user *model.Comment) error
	DeleteComment(ctx context.Context, comment *model.Comment) (gen.ResultInfo, error)
}

type CommentRepository struct {
	CommentQuery *comment
}

func (c *CommentRepository) GetCommentById(ctx context.Context, userId int64) (*model.Comment, error) {
	return c.CommentQuery.WithContext(ctx).Where(c.CommentQuery.ID.Eq(userId)).First()
}

func (c *CommentRepository) GetCommentListByVideoId(ctx context.Context, videoId int64) ([]*model.Comment, error) {
	return c.CommentQuery.WithContext(ctx).Where(c.CommentQuery.VideoID.Eq(videoId)).Find()
}

func (c *CommentRepository) GetCommentCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	return c.CommentQuery.WithContext(ctx).Where(c.CommentQuery.VideoID.Eq(videoId)).Count()
}

func (c *CommentRepository) InsertComment(ctx context.Context, user *model.Comment) error {
	return c.CommentQuery.WithContext(ctx).Create(user)
}

func (c *CommentRepository) DeleteComment(ctx context.Context, comment *model.Comment) (gen.ResultInfo, error) {
	return c.CommentQuery.WithContext(ctx).Delete(comment)
}

func NewCommentRepository(dsn string) CommentDo {
	q := initCommentQuery(dsn)
	return &CommentRepository{
		CommentQuery: q,
	}
}

// 初始化数据库设置
func initCommentQuery(dsn string) *comment {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 跳过默认事务
		SkipDefaultTransaction: false,
		// 命名策略（表名、列名生成规则）
		NamingStrategy: nil,
		// 创建更新时候，是否更新关联数据
		FullSaveAssociations: false,
		// 创建时间函数
		NowFunc: nil,
		// 生成SQL不执行
		DryRun: false,
		// 是否禁止创建prepare stm
		PrepareStmt: false,
		// 禁用数据库健康检查
		DisableAutomaticPing: true,
		// 是否禁止自动创建外间约束
		DisableForeignKeyConstraintWhenMigrating: false,
		// 是否自动禁止外键约束
		IgnoreRelationshipsWhenMigrating: false,
		// 是否禁止嵌套事务
		DisableNestedTransaction: false,
		// 是否允许全局更新
		AllowGlobalUpdate: false,
		// 查询是否带上全部字段
		QueryFields: false,
		// 默认批量插入大小
		CreateBatchSize: 400,
		ClauseBuilders:  nil,
		ConnPool:        nil,
		Dialector:       nil,
		Plugins:         nil,
	})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Minute * 20)
	// 注入dao
	q := Use(db)
	return &q.Comment
}
