package dao

import (
	"context"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/common/xconst"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"time"
)

type VideoDo interface {
	GetVideoById(ctx context.Context, videoId int64) (*model.Video, error)
	GetVideoListByUserId(ctx context.Context, userId int64) ([]*model.Video, error)
	GetAllVideo(ctx context.Context) ([]*model.Video, error)
	GetVideoListByTimeStampAndSectionId(ctx context.Context, timestamp int64, sectionId int64) ([]*model.Video, error)
	InsertVideo(ctx context.Context, video *model.Video) error
	DeleteVideo(ctx context.Context, video *model.Video) (gen.ResultInfo, error)
}

type VideoRepository struct {
	VideoQuery *video
}

func (v *VideoRepository) GetVideoById(ctx context.Context, videoId int64) (*model.Video, error) {
	return v.VideoQuery.WithContext(ctx).Where(v.VideoQuery.ID.Eq(videoId)).First()
}

func (v *VideoRepository) GetVideoListByUserId(ctx context.Context, userId int64) ([]*model.Video, error) {
	return v.VideoQuery.WithContext(ctx).Where(v.VideoQuery.OwnerID.Eq(userId)).Find()
}

func (v *VideoRepository) GetAllVideo(ctx context.Context) ([]*model.Video, error) {
	return v.VideoQuery.WithContext(ctx).Find()
}

func (v *VideoRepository) InsertVideo(ctx context.Context, video *model.Video) error {
	return v.VideoQuery.WithContext(ctx).Create(video)
}

func (v *VideoRepository) DeleteVideo(ctx context.Context, video *model.Video) (gen.ResultInfo, error) {
	return v.VideoQuery.WithContext(ctx).Delete(video)
}

func (v *VideoRepository) GetVideoListByTimeStampAndSectionId(ctx context.Context, timestamp int64, sectionId int64) ([]*model.Video, error) {
	return v.VideoQuery.WithContext(ctx).Where(v.VideoQuery.CreateTime.Lt(timestamp)).Where(v.VideoQuery.SectionID.Eq(sectionId)).Limit(xconst.VideoFeedCount).Order(v.VideoQuery.CreateTime.Desc()).Find()
}

func NewVideoRepository(dsn string) VideoDo {
	q := initVideoQuery(dsn)
	return &VideoRepository{
		VideoQuery: q,
	}
}

// 初始化数据库设置
func initVideoQuery(dsn string) *video {
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
	return &q.Video
}
