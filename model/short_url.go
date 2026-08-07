package model

import (
	"reflect"
	"time"
)

// ShortUrl 短链接
type ShortUrl struct {
	ID          int64     `db:"id" json:"id"`
	LongUrl     string    `db:"long_url" json:"long_url"`   //原始链接
	ShortUrl    string    `db:"short_url" json:"short_url"` //生成的短链接
	CreatedTime time.Time `db:"created_time" json:"created_time"`
	Comment     string    `db:"comment" json:"comment"` //备注
}

func (url ShortUrl) IsEmpty() bool {
	return reflect.DeepEqual(url, ShortUrl{})
}
