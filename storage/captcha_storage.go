package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/dchest/captcha"
)

type CaptchaRedisStore struct {
	RedisService *RedisService
	Expiration   time.Duration
	KeyPrefix    string
}

func NewRedisStore(rs *RedisService, expiration time.Duration, prefixKey string) captcha.Store {
	s := &CaptchaRedisStore{
		RedisService: rs,
		Expiration:   expiration,
		KeyPrefix:    prefixKey,
	}
	return s
}

func (s *CaptchaRedisStore) Set(id string, digit []byte) {
	key := fmt.Sprintf("%v#%v", s.KeyPrefix, id)
	ctx := context.Background()

	err := RedisSet(ctx, key, digit, s.Expiration)
	if err != nil {
		panic("redis set error")
	}
}

func (s *CaptchaRedisStore) Get(id string, clear bool) []byte {
	key := fmt.Sprintf("%v#%v", s.KeyPrefix, id)
	ctx := context.Background()

	result, err := RedisGetString(ctx, key)
	if err != nil {
		return nil
	}
	if clear {
		RedisDelete(ctx, key)
	}
	return []byte(result)
}
