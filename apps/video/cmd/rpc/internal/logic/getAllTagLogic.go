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

type GetAllTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllTagLogic {
	return &GetAllTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllTagLogic) GetAllTag(in *pb.GetAllTagReq) (*pb.GetAllTagResp, error) {
	// todo: add your logic here and delete this line
	// 直接从redis中获取数据
	exists, err := l.svcCtx.Redis.Exists(l.ctx, xconst.RedisVideoTag)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("判断redis视频标签的键是否存在时失败，err: %v", err)
	}

	// 如果redis中有数据，直接返回
	if exists {
		if tags, err := l.svcCtx.Redis.Smembers(l.ctx, xconst.RedisVideoTag); err != nil {
			l.Logger.Errorf("从redis中获取所有视频标签失败，err: %v", err)
		} else {
			if len(tags) == 0 {
				return &pb.GetAllTagResp{}, nil
			}
			resp := &pb.GetAllTagResp{Tags: make([]*pb.TagInfo, 0, len(tags))}
			for _, tag := range tags {
				tagModel := &model.Tag{}
				if err := json.Unmarshal([]byte(tag), tagModel); err != nil {
					l.Logger.Errorf("反序列化redis标签信息时失败, err: %v", err)
					continue
				}
				resp.Tags = append(resp.Tags, &pb.TagInfo{Id: tagModel.ID, Name: tagModel.Name})
			}
			// 更新缓存失效时间
			err := l.svcCtx.Redis.Expire(l.ctx, xconst.RedisVideoTag, xconst.RedisExpireTime)
			if err != nil {
				l.Logger.Errorf("设置redis分区缓存失效时间失败, err: %v", err)
			}
			return resp, nil
		}
	}

	// 从数据库中获取数据
	key := utils.GetRedisKeyWithPrefix(xconst.RedisVideoTag, 1)
	tags, err := l.svcCtx.SingleFlight.Do(key, func() (any, error) {
		return l.svcCtx.TagDo.GetAllTags(l.ctx)
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "从数据库中获取所有视频标签信息失败, err: %v", err)
	}

	// 类型断言
	tagsModel, ok := tags.([]*model.Tag)
	if !ok {
		return nil, errors.Wrapf(xerr.NewErrMsg("所有视频标签的类型断言失败"), "data: %v", tags)
	}

	// 长度为零
	if len(tagsModel) == 0 {
		return &pb.GetAllTagResp{}, nil
	}

	// 拼接数据
	resp := &pb.GetAllTagResp{Tags: make([]*pb.TagInfo, 0, len(tagsModel))}
	for _, tag := range tagsModel {
		tagInfo := &pb.TagInfo{}
		tagInfo.Id = tag.ID
		tagInfo.Name = tag.Name
		resp.Tags = append(resp.Tags, tagInfo)
	}

	tsg := utils.NewTestGo()
	tsg.RunSafe(l.BuildVideoTagCache)

	return resp, nil
}

func (l *GetAllTagLogic) BuildVideoTagCache() {
	// 获取分布式锁键
	lock := l.svcCtx.Redis.NewRedisLock(xconst.RedisBuildVideoTagCacheLockKey)
	lock.SetExpire(1)

	// 复制ctx，防止异步调用时logic.ctx结束
	ctx := contextx.ValueOnlyFrom(l.ctx)

	// 获取分布式锁
	acquire, err := lock.Acquire()
	if err != nil {
		logx.WithContext(ctx).Errorf("获取分布式锁失败，lockKey: %s, err: %v", xconst.RedisBuildVideoTagCacheLockKey, err)
		return
	}

	// 延迟释放分布式锁
	defer func(lock *redis.RedisLock) {
		_, err := lock.Release()
		if err != nil {
			logx.WithContext(ctx).Errorf("释放分布式锁失败，lockKey: %s, err: %v", xconst.RedisBuildVideoTagCacheLockKey, err)
		}
	}(lock)

	// 更新缓存
	if acquire {
		// 从数据库中查询所有标签信息
		tags, err := l.svcCtx.TagDo.GetAllTags(ctx)
		if err != nil {
			logx.WithContext(ctx).Errorf("get all tags failed, err: %v", err)
		}
		if len(tags) > 0 {
			tagList := make([]interface{}, 0, len(tags))
			for _, tag := range tags {
				if tagString, err := json.Marshal(tag); err != nil {
					logx.WithContext(ctx).Errorf("序列化redis标签信息时失败, err: %v", err)
				} else {
					tagList = append(tagList, tagString)
				}
			}

			// 设置缓存
			_, err := l.svcCtx.Redis.Sadd(ctx, xconst.RedisVideoTag, tagList...)
			if err != nil {
				logx.WithContext(ctx).Errorf("add redis video tag cache  failed, err: %v", err)
				return
			}

			// 设置缓存失效时间
			err = l.svcCtx.Redis.Expire(ctx, xconst.RedisVideoTag, xconst.RedisExpireTime)
			if err != nil {
				logx.WithContext(ctx).Errorf("set video tag redis key expire time failed, err: %v", err)
				return
			}
		}
		return
	}
}
