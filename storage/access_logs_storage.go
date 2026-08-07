package storage

import (
	"context"
	"my_url_shortner/model"
)

// GetAccessLogs 获取所有日志信息
func GetAccessLogs(ctx context.Context) ([]model.AccessLog, error) {
	query := `SELECT * FROM public.access_logs`

	var res []model.AccessLog

	return res, DbSelect(ctx, query, &res)
}

func GetAccessLogsByPage(ctx context.Context, page, pageSize int) ([]model.AccessLog, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	query := `SELECT * FROM public.access_logs ORDER BY id ASC LIMIT $1 OFFSET $2`

	var res []model.AccessLog

	return res, DbSelect(ctx, query, &res, pageSize, (page-1)*pageSize)
}
