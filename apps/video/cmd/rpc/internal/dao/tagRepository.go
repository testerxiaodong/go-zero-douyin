package dao

import (
	"context"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"time"
)

type TagDo interface {
	GetTagById(ctx context.Context, tagId int64) (*model.Tag, error)
	GetTagByName(ctx context.Context, name string) (*model.Tag, error)
	GetAllTags(ctx context.Context) ([]*model.Tag, error)
	InsertTag(ctx context.Context, tag *model.Tag) error
	DeleteTag(ctx context.Context, tag *model.Tag) (gen.ResultInfo, error)
}

type TagRepository struct {
	TagQuery *tag
}

func (s *TagRepository) GetTagById(ctx context.Context, tagId int64) (*model.Tag, error) {
	return s.TagQuery.WithContext(ctx).Where(s.TagQuery.ID.Eq(tagId)).First()
}

func (s *TagRepository) GetTagByName(ctx context.Context, name string) (*model.Tag, error) {
	return s.TagQuery.WithContext(ctx).Where(s.TagQuery.Name.Eq(name)).First()
}

func (s *TagRepository) GetAllTags(ctx context.Context) ([]*model.Tag, error) {
	return s.TagQuery.WithContext(ctx).Find()
}

func (s *TagRepository) InsertTag(ctx context.Context, tag *model.Tag) error {
	return s.TagQuery.WithContext(ctx).Create(tag)
}

func (s *TagRepository) DeleteTag(ctx context.Context, tag *model.Tag) (gen.ResultInfo, error) {
	return s.TagQuery.WithContext(ctx).Delete(tag)
}

func NewTagRepository(dsn string) TagDo {
	q := initTagQuery(dsn)
	return &TagRepository{
		TagQuery: q,
	}
}

// 初始化数据库设置
func initTagQuery(dsn string) *tag {
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
	return &q.Tag
}
