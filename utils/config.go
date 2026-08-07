package utils

import (
	"my_url_shortner/global"

	"gopkg.in/ini.v1"
)

func InitConfig(file string) error {
	cfg, err := ini.Load(file)
	if err != nil {
		return err
	}

	postgreSection := cfg.Section("postgres")
	global.Conf.DataBase.Host = postgreSection.Key("host").String()
	global.Conf.DataBase.Port = postgreSection.Key("port").MustInt()
	global.Conf.DataBase.User = postgreSection.Key("user").String()
	global.Conf.DataBase.Password = postgreSection.Key("password").String()
	global.Conf.DataBase.DbName = postgreSection.Key("dbname").String()
	global.Conf.DataBase.MaxOpenConns = postgreSection.Key("max_open_conns").MustInt()
	global.Conf.DataBase.MaxIdleTime = postgreSection.Key("max_idle_time").MustInt()

	redisSection := cfg.Section("redis")
	global.Conf.Redis.Host = redisSection.Key("host").String()
	global.Conf.Redis.Database = redisSection.Key("database").MustInt()
	global.Conf.Redis.User = redisSection.Key("username").String()
	global.Conf.Redis.Password = redisSection.Key("password").String()
	global.Conf.Redis.PoolSize = redisSection.Key("pool_size").MustInt()

	appSection := cfg.Section("app")
	global.Conf.App.Port = appSection.Key("port").MustInt()
	global.Conf.App.AdminPort = appSection.Key("adminport").MustInt()
	global.Conf.App.UrlPrefix = appSection.Key("url_prefix").String()

	captchaSection := cfg.Section("captcha")
	global.Conf.Captcha.Strore = captchaSection.Key("store").String()

	return nil
}
