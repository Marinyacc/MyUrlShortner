package model

import (
	"time"
)

// AccessLog 访问日志
type AccessLog struct {
	ID          int64     `db:"id"`
	ShortUrl    string    `db:"short_url"`
	CreatedTime time.Time `db:"created_time"`
	Ip          string    `db:"ip"`
}
