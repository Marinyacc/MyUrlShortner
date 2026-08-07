package global

import (
	"my_url_shortner/model"
	"time"
)

const (
	ShortUrl_CountPrefix    = "Url_Count:"   //某个短链接总点击量
	ShortUrl_IP_CountPrefix = "Url_d_Count:" //某个短链接总不同IP点击量
	ShortUrl_ID             = "Url_ID:"      // 自增ID前缀
	ShortUrl_Active_Set     = "active_set:"
	ShortUrl                = "Url:"      //[Url:shortUrl,longUrl]
	ShortUrl_Lock           = "Url_Lock:" //某个短链接的分布式锁
	PvTotalKey              = "pvTotalKey"
	UvTotalKey              = "uvTotalKey"
	UpdateTime              = 1 * time.Minute //统计数据定时任务redis -> postgres 时间
	UpdateTime2             = 1 * time.Minute //统计数据定时任务总表
	CookieTTL               = 60 * 30         //token 30分钟过期
	StreamName              = "access_logs_stream"
	GroupName               = "access_logs_group"
	ConsumerName            = "server"
	Page                    = 1
	PageSize                = 10
)

var (
	Conf model.Config
)
