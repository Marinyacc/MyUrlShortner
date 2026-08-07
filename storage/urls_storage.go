package storage

import (
	"context"
	"my_url_shortner/global"
	"my_url_shortner/model"

	"github.com/jackc/pgx/v5"
)


// GetUrlInfo 获取指定url的信息
func GetUrlInfo(ctx context.Context, shorturl string) (model.ShortUrl, error) {
	var resp model.ShortUrl
	query := `SELECT * FROM public.short_urls WHERE short_url = $1`
	return resp, DbGet(ctx, query, &resp, shorturl)
}

// GetUrlsInfo 获取所有url的信息
func GetUrlsInfo(ctx context.Context) ([]model.ShortUrl, error) {
	var resp []model.ShortUrl
	query := `SELECT * FROM public.short_urls`
	return resp, DbSelect(ctx, query, &resp)
}

func GetUrlsInfoByPage(ctx context.Context, page, pageSize int) ([]model.ShortUrl, error) {
	var resp []model.ShortUrl
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	query := `SELECT * FROM public.short_urls LIMIT $1 OFFSET $2`
	return resp, DbSelect(ctx, query, &resp, pageSize, (page-1)*pageSize)
}

// GetUrlStats 获取指定url的统计信息
func GetUrlStats(ctx context.Context, shorturl string) (model.ShortUrlStats, error) {
	var resp model.ShortUrlStats
	query := `SELECT * FROM public.stats WHERE short_url = $1`
	return resp, DbGet(ctx, query, &resp, shorturl)
}

// GetUrlsStats 获取所有url的统计信息
func GetUrlsStats(ctx context.Context) ([]model.ShortUrlStats, error) {
	var resp []model.ShortUrlStats
	query := `SELECT * FROM public.stats`
	return resp, DbSelect(ctx, query, &resp)
}

func GetUrlsStatsByPage(ctx context.Context, page, pageSize int) ([]model.ShortUrlStats, error) {
	var resp []model.ShortUrlStats
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	query := `SELECT * FROM public.stats LIMIT $1 OFFSET $2 `
	return resp, DbSelect(ctx, query, &resp, pageSize, (page-1)*pageSize)
}

// InsertShortUrl 插入短链接
func InsertShortUrl(ctx context.Context, shortUrl model.ShortUrl) error {
	query := `INSERT INTO public.short_urls (short_url, long_url,comment)
	 VALUES(@short_url,@long_url,@comment)`
	return DbNamedExec(ctx, query, pgx.NamedArgs{
		"short_url": shortUrl.ShortUrl,
		"long_url":  shortUrl.LongUrl,
		"comment":   shortUrl.Comment,
	})
}

// GenerateShortUrl 采用自增ID + base62 编码生成短链接
func GenerateShortUrl(ctx context.Context, LongUrl string) (string, error) {
	var (
		ID  int64
		err error
	)
	ID, err = Rs.redisClient.Incr(ctx, global.ShortUrl_ID).Result()

	if err != nil {
		return "", err
	}

	return encode62(ID), nil
}

func encode62(id int64) string {
	const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if id == 0 {
		return string(base62Chars[0])
	}
	var res []byte

	for id > 0 {
		idx := id % 62
		res = append(res, base62Chars[idx])
		id /= 62
	}

	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}

	return string(res)
}
