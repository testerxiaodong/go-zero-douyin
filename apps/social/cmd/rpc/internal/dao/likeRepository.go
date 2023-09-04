package dao

import (
	"context"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"time"
)

type LikeDo interface {
	GetLikeByVideoIdAndUserId(ctx context.Context, videoId int64, userId int64) (*model.Like, error)
	GetVideoLikedCount(ctx context.Context, videoId int64) (int64, error)
	GetUserLikeVideoIdList(ctx context.Context, UserId int64) ([]int64, error)
	GetVideoLikedByUserIdList(ctx context.Context, videoId int64) ([]int64, error)
	InsertLike(ctx context.Context, user *model.Like) error
	DeleteLike(ctx context.Context, like *model.Like) (gen.ResultInfo, error)
}

type LikeRepository struct {
	LikeQuery *like
}

func (f *LikeRepository) GetLikeByVideoIdAndUserId(ctx context.Context, videoId int64, userId int64) (*model.Like, error) {
	return f.LikeQuery.WithContext(ctx).Where(f.LikeQuery.VideoID.Eq(videoId)).Where(f.LikeQuery.UserID.Eq(userId)).First()
}

func (f *LikeRepository) GetVideoLikedCount(ctx context.Context, videoId int64) (int64, error) {
	return f.LikeQuery.WithContext(ctx).Where(f.LikeQuery.VideoID.Eq(videoId)).Count()
}

func (f *LikeRepository) GetUserLikeVideoIdList(ctx context.Context, userId int64) ([]int64, error) {
	likes, err := f.LikeQuery.WithContext(ctx).Where(f.LikeQuery.UserID.Eq(userId)).Find()
	if err != nil {
		return nil, err
	}
	if len(likes) > 0 {
		idList := make([]int64, 0)
		for _, like := range likes {
			idList = append(idList, like.VideoID)
		}
		return idList, nil
	}
	return []int64{}, nil
}

func (f *LikeRepository) GetVideoLikedByUserIdList(ctx context.Context, videoId int64) ([]int64, error) {
	likes, err := f.LikeQuery.WithContext(ctx).Where(f.LikeQuery.VideoID.Eq(videoId)).Find()
	if err != nil {
		return nil, err
	}
	if len(likes) > 0 {
		idList := make([]int64, 0)
		for _, like := range likes {
			idList = append(idList, like.UserID)
		}
		return idList, nil
	}
	return []int64{}, nil
}

func (f *LikeRepository) InsertLike(ctx context.Context, user *model.Like) error {
	return f.LikeQuery.WithContext(ctx).Create(user)
}

func (f *LikeRepository) DeleteLike(ctx context.Context, like *model.Like) (gen.ResultInfo, error) {
	return f.LikeQuery.WithContext(ctx).Delete(like)
}

func NewLikeRepository(dsn string) LikeDo {
	q := initLikeQuery(dsn)
	return &LikeRepository{
		LikeQuery: q,
	}
}

// 初始化数据库设置
func initLikeQuery(dsn string) *like {
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
	return &q.Like
}
