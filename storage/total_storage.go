package storage

import (
	"context"
)

func TotalUrls(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM public.short_urls`
	var resp int
	return resp, DbGet(ctx, query, &resp)
}

func TotalStats(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM public.stats`
	var resp int
	return resp, DbGet(ctx, query, &resp)
}

func TotalAccessLogs(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM public.access_logs`
	var resp int
	return resp, DbGet(ctx, query, &resp)
}
