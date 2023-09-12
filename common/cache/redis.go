package cache

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type RedisCache interface {
	Exists(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, seconds int) error
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) (int, error)
	Scard(ctx context.Context, key string) (int64, error)
	Sadd(ctx context.Context, key string, values ...any) (int, error)
	Smembers(ctx context.Context, key string) ([]string, error)
	NewRedisLock(key string) *redis.RedisLock
}

type RedisLocker interface {
	SetExpire(seconds int)
	Acquire() (bool, error)
	AcquireCtx(ctx context.Context) (bool, error)
	Release() (bool, error)
	ReleaseCtx(ctx context.Context) (bool, error)
}

type RedisLock struct {
	lock *redis.RedisLock
}

func NewRedisLock(lock *redis.RedisLock) RedisLocker {
	return &RedisLock{
		lock: lock,
	}
}

func (r *RedisLock) SetExpire(seconds int) {
	r.lock.SetExpire(seconds)
}

func (r *RedisLock) Acquire() (bool, error) {
	return r.lock.Acquire()
}

func (r *RedisLock) AcquireCtx(ctx context.Context) (bool, error) {
	return r.lock.AcquireCtx(ctx)
}

func (r *RedisLock) Release() (bool, error) {
	return r.lock.Release()
}

func (r *RedisLock) ReleaseCtx(ctx context.Context) (bool, error) {
	return r.lock.ReleaseCtx(ctx)
}

type RedisClient struct {
	Redis *redis.Redis
}

func NewRedisClient(c redis.RedisConf) *RedisClient {
	return &RedisClient{
		Redis: redis.MustNewRedis(c),
	}
}

func (rc *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	return rc.Redis.ExistsCtx(ctx, key)
}

func (rc *RedisClient) Expire(ctx context.Context, key string, seconds int) error {
	return rc.Redis.ExpireCtx(ctx, key, seconds)
}

func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return rc.Redis.GetCtx(ctx, key)
}

func (rc *RedisClient) Set(ctx context.Context, key string, value string) error {
	return rc.Redis.SetCtx(ctx, key, value)
}

func (rc *RedisClient) Delete(ctx context.Context, key string) (int, error) {
	return rc.Redis.DelCtx(ctx, key)
}

func (rc *RedisClient) Scard(ctx context.Context, key string) (int64, error) {
	return rc.Redis.ScardCtx(ctx, key)
}

func (rc *RedisClient) Sadd(ctx context.Context, key string, values ...any) (int, error) {
	return rc.Redis.SaddCtx(ctx, key, values)
}

func (rc *RedisClient) Smembers(ctx context.Context, key string) ([]string, error) {
	return rc.Redis.SmembersCtx(ctx, key)
}

func (rc *RedisClient) NewRedisLock(key string) *redis.RedisLock {
	return redis.NewRedisLock(rc.Redis, key)
}
