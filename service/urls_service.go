package service

import (
	"context"
	"errors"
	"my_url_shortner/global"
	"my_url_shortner/model"
	"my_url_shortner/storage"
)

// GenerateShortUrl 生成短链接,并插入到Redis和数据库中
func GenerateShortUrl(ctx context.Context, longUrl, comment string) (string, error) {
	shortUrl, err := storage.GenerateShortUrl(ctx, longUrl)
	if err != nil {
		return "", err
	}

	//1.Redis 键值对 Key:[Url:shortUrl] Value:[longUrl]
	if err := storage.RedisSet4Ever(ctx, global.ShortUrl+shortUrl, longUrl); err != nil {
		return "", errors.New("Redis 设置键值对出错")
	}

	//2. Postgresql
	URL := model.ShortUrl{
		ShortUrl: shortUrl,
		LongUrl:  longUrl,
		Comment:  comment,
	}
	if err := storage.InsertShortUrl(ctx, URL); err != nil {
		return "", err
	}
	return shortUrl, nil
}
