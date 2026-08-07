package storage

import (
	"context"
	"log"
	"my_url_shortner/utils"
	"time"
)

const script = `
	if redis.call("get",KEYS[1]) == ARGV[1] then
		return redis.call("del",KEYS[1])
	else 
		return 0
	end
`

// Lock redis分布式锁结构
type Lock struct {
	key   string
	value string
	ttl   time.Duration
}

// NewLock 获取锁的新实例
func NewLock(ctx context.Context, key string, ttl time.Duration) *Lock {
	return &Lock{
		key:   key,
		value: utils.RandString(16),
		ttl:   ttl,
	}
}

// 尝试上锁
func (l *Lock) Lock(ctx context.Context) bool {
	var (
		result bool
		err    error
	)
	result, err = Rs.redisClient.SetNX(ctx, l.key, l.value, l.ttl).Result()
	if err != nil {
		log.Println("Redis 分布式锁上锁出错!:", err.Error())
		return false
	}
	return result
}

// 尝试解锁
func (l *Lock) UnLock(ctx context.Context) bool {
	var (
		result any
		err    error
	)
	result, err = Rs.redisClient.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		log.Println("Redis 分布式锁解锁出错!:", err.Error())
		return false
	}
	val, _ := result.(int64)
	return val == 1
}
