package logic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/contextx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllSectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllSectionLogic {
	return &GetAllSectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllSectionLogic) GetAllSection(in *pb.GetAllSectionReq) (*pb.GetAllSectionResp, error) {
	// todo: add your logic here and delete this line
	// 直接从redis中获取数据
	exists, err := l.svcCtx.Redis.Exists(l.ctx, xconst.RedisVideoSection)
	if err != nil {
		l.Logger.Errorf("判断redis视频分区的键是否存在时失败，err: %v", err)
	}

	// 如果redis中有数据，直接返回
	if exists {
		if sections, err := l.svcCtx.Redis.Smembers(l.ctx, xconst.RedisVideoSection); err != nil {
			l.Logger.Errorf("从redis中获取所有视频分区失败，err: %v", err)
		} else {
			if len(sections) == 0 {
				return &pb.GetAllSectionResp{}, nil
			}
			resp := &pb.GetAllSectionResp{Sections: make([]*pb.SectionInfo, 0, len(sections))}
			for _, section := range sections {
				sectionModel := &model.Section{}
				if err := json.Unmarshal([]byte(section), sectionModel); err != nil {
					l.Logger.Errorf("反序列化redis分区信息时失败, err: %v", err)
					continue
				}
				resp.Sections = append(resp.Sections, &pb.SectionInfo{Id: sectionModel.ID, Name: sectionModel.Name})
			}
			err := l.svcCtx.Redis.Expire(l.ctx, xconst.RedisVideoSection, xconst.RedisExpireTime)
			if err != nil {
				l.Logger.Errorf("设置redis分区缓存失效时间失败, err: %v", err)
			}
			return resp, nil
		}
	}

	// 从数据库中获取数据
	key := utils.GetRedisKeyWithPrefix(xconst.RedisVideoSection, 1)
	sections, err := l.svcCtx.SingleFlight.Do(key, func() (any, error) {
		return l.svcCtx.SectionDo.GetAllSections(l.ctx)
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "从数据库中获取所有视频分区信息失败, err: %v", err)
	}

	// 类型断言
	sectionsModel, ok := sections.([]*model.Section)
	if !ok {
		return nil, errors.Wrapf(xerr.NewErrMsg("所有视频分区的类型断言失败"), "data: %v", sections)
	}

	// 长度为零
	if len(sectionsModel) == 0 {
		return &pb.GetAllSectionResp{}, nil
	}

	// 拼接数据
	resp := &pb.GetAllSectionResp{Sections: make([]*pb.SectionInfo, 0, len(sectionsModel))}
	for _, section := range sectionsModel {
		sectionInfo := &pb.SectionInfo{}
		sectionInfo.Id = section.ID
		sectionInfo.Name = section.Name
		resp.Sections = append(resp.Sections, sectionInfo)
	}

	tsg := utils.NewTestGo()
	tsg.RunSafe(l.BuildVideoSectionCache)

	return resp, nil
}

func (l *GetAllSectionLogic) BuildVideoSectionCache() {
	// 获取分布式锁键
	lock := l.svcCtx.Redis.NewRedisLock(xconst.RedisBuildVideoSectionCacheLockKey)
	lock.SetExpire(1)

	// 复制ctx，防止异步调用时logic.ctx结束
	ctx := contextx.ValueOnlyFrom(l.ctx)

	// 获取分布式锁
	acquire, err := lock.Acquire()
	if err != nil {
		logx.WithContext(ctx).Errorf("获取分布式锁失败，lockKey: %s, err: %v", xconst.RedisBuildVideoSectionCacheLockKey, err)
		return
	}

	// 延迟释放分布式锁
	defer func(lock *redis.RedisLock) {
		_, err := lock.Release()
		if err != nil {
			logx.WithContext(ctx).Errorf("释放分布式锁失败，lockKey: %s, err: %v", xconst.RedisBuildVideoSectionCacheLockKey, err)
		}
	}(lock)

	// 更新缓存
	if acquire {
		// 从数据库中查询视频点赞用户列表
		sections, err := l.svcCtx.SectionDo.GetAllSections(ctx)
		if err != nil {
			logx.WithContext(ctx).Errorf("get all sections failed, err: %v", err)
		}
		if len(sections) > 0 {
			sectionList := make([]interface{}, 0, len(sections))
			for _, section := range sections {
				if sectionString, err := json.Marshal(section); err != nil {
					logx.WithContext(ctx).Errorf("序列化redis分区信息时失败, err: %v", err)
				} else {
					sectionList = append(sectionList, sectionString)
				}
			}

			// 设置缓存
			_, err := l.svcCtx.Redis.Sadd(ctx, xconst.RedisVideoSection, sectionList...)
			if err != nil {
				logx.WithContext(ctx).Errorf("add redis video section cache  failed, err: %v", err)
				return
			}

			// 设置缓存失效时间
			err = l.svcCtx.Redis.Expire(ctx, xconst.RedisVideoSection, xconst.RedisExpireTime)
			if err != nil {
				logx.WithContext(ctx).Errorf("set video section redis key expire time failed, err: %v", err)
				return
			}
		}
		return
	}
}
