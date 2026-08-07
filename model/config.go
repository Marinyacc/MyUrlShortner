package model

type Config struct {
	DataBase DataBaseConfig
	Captcha  CaptchaConfig
	App      AppConfig
	Redis    RedisConfig
}

type CaptchaConfig struct {
	Strore string
}

type AppConfig struct {
	Port      int
	AdminPort int
	UrlPrefix string
}

type RedisConfig struct {
	Host     string
	User     string
	Password string
	Database int
	PoolSize int
}

type DataBaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DbName       string
	MaxOpenConns int
	MaxIdleTime  int
}
