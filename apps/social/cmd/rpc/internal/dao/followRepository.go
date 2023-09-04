package dao

import (
	"context"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"time"
)

type FollowDo interface {
	GetFollowByFollowerIdAndUserId(ctx context.Context, followerId int64, userId int64) (*model.Follow, error)
	GetUserFollowCount(ctx context.Context, videoId int64) (int64, error)
	GetUserFollowerCount(ctx context.Context, videoId int64) (int64, error)
	GetUserFollowIdList(ctx context.Context, followerId int64) ([]int64, error)
	GetUserFollowerIdList(ctx context.Context, userId int64) ([]int64, error)
	InsertFollow(ctx context.Context, user *model.Follow) error
	DeleteFollow(ctx context.Context, follow *model.Follow) (gen.ResultInfo, error)
}

type FollowRepository struct {
	FollowQuery *follow
}

func (f *FollowRepository) GetFollowByFollowerIdAndUserId(ctx context.Context, followerId int64, userId int64) (*model.Follow, error) {
	return f.FollowQuery.WithContext(ctx).Where(f.FollowQuery.FollowerID.Eq(followerId)).Where(f.FollowQuery.UserID.Eq(userId)).First()
}

func (f *FollowRepository) GetUserFollowCount(ctx context.Context, followerId int64) (int64, error) {
	return f.FollowQuery.WithContext(ctx).Where(f.FollowQuery.FollowerID.Eq(followerId)).Count()
}

func (f *FollowRepository) GetUserFollowerCount(ctx context.Context, userId int64) (int64, error) {
	return f.FollowQuery.WithContext(ctx).Where(f.FollowQuery.UserID.Eq(userId)).Count()
}

func (f *FollowRepository) GetUserFollowIdList(ctx context.Context, followerId int64) ([]int64, error) {
	follows, err := f.FollowQuery.WithContext(ctx).Where(f.FollowQuery.FollowerID.Eq(followerId)).Find()
	if err != nil {
		return nil, err
	}
	if len(follows) > 0 {
		idList := make([]int64, 0)
		for _, follow := range follows {
			idList = append(idList, follow.UserID)
		}
		return idList, nil
	}
	return []int64{}, nil
}

func (f *FollowRepository) GetUserFollowerIdList(ctx context.Context, userId int64) ([]int64, error) {
	follows, err := f.FollowQuery.WithContext(ctx).Where(f.FollowQuery.UserID.Eq(userId)).Find()
	if err != nil {
		return nil, err
	}
	if len(follows) > 0 {
		idList := make([]int64, 0)
		for _, follow := range follows {
			idList = append(idList, follow.FollowerID)
		}
		return idList, nil
	}
	return []int64{}, nil
}

func (f *FollowRepository) InsertFollow(ctx context.Context, user *model.Follow) error {
	return f.FollowQuery.WithContext(ctx).Create(user)
}

func (f *FollowRepository) DeleteFollow(ctx context.Context, follow *model.Follow) (gen.ResultInfo, error) {
	return f.FollowQuery.WithContext(ctx).Delete(follow)
}

func NewFollowRepository(dsn string) FollowDo {
	q := initFollowQuery(dsn)
	return &FollowRepository{
		FollowQuery: q,
	}
}

// 初始化数据库设置
func initFollowQuery(dsn string) *follow {
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
	return &q.Follow
}
