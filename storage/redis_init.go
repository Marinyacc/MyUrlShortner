package storage

import (
	"context"
	"my_url_shortner/global"
	"time"

	"github.com/redis/go-redis/v9"
)

// redis 服务
var Rs *RedisService = &RedisService{}

type RedisService struct {
	redisClient *redis.Client
}

// 初始化 Redis服务
func InitRedisServcie() (*RedisService, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     global.Conf.Redis.Host,
		Username: global.Conf.Redis.User,
		Password: global.Conf.Redis.Password,
		DB:       global.Conf.Redis.Database, //Redis 默认有16个数据库(0~15),默认选择0号数据库.不同数据库数据相互独立，若使用集群模式，则数据库分类被禁用
		PoolSize: global.Conf.Redis.PoolSize,
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	Rs.redisClient = redisClient

	return Rs, nil
}

// RedisSet 设置 Redis键值对
func RedisSet(ctx context.Context, key string, value any, ttl time.Duration) error {
	return Rs.redisClient.Set(ctx, key, value, ttl).Err()
}

// RedisSet30m 设置Redis键值对，过期时间为30分钟
func RedisSet30m(ctx context.Context, key string, value any) error {
	return RedisSet(ctx, key, value, 30*time.Minute)
}

// RedisSet4Ever 设置Redis键值对，永不过期
func RedisSet4Ever(ctx context.Context, key string, value any) error {
	return RedisSet(ctx, key, value, 0)
}

// RedisGetString 获取 Redis字符串
func RedisGetString(ctx context.Context, key string) (string, error) {
	var (
		result string
		err    error
	)
	result, err = Rs.redisClient.Get(ctx, key).Result()

	if err == redis.Nil {
		return result, nil
	}
	return result, err
}

// RedisDelete 删除Redis中的键值对
func RedisDelete(ctx context.Context, key ...string) error {
	if len(key) > 0 {
		return Rs.redisClient.Del(ctx, key...).Err()
	}
	return nil
}

// 关闭RedisService
func RedisClose() {
	Rs.redisClient.Close()
}


