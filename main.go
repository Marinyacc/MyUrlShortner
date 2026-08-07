package main

import (
	"context"
	"fmt"
	"log"
	"my_url_shortner/global"
	"my_url_shortner/router"
	"my_url_shortner/storage"
	"my_url_shortner/utils"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

var ctx context.Context = context.Background()

func main() {

	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	if err := utils.InitConfig("config.ini"); err != nil {
		panic(err)
	}
	log.Println("成功解析配置文件")

	rs, err := storage.InitRedisServcie()
	if err != nil {
		panic(err)
	}
	log.Println("redis服务成功启动")
	defer storage.RedisClose()

	pool, err := storage.InitPostgreService()
	if err != nil {
		panic(err)
	}
	log.Println("postgre服务成功启动")
	defer storage.DbClose()

	if strings.EqualFold("redis", strings.ToLower(global.Conf.Captcha.Strore)) {
		crs := storage.NewRedisStore(rs, 1*time.Minute, "my_captcha")
		captcha.SetCustomStore(crs)
	}
	log.Println("captcha服务成功启动")

	var group errgroup.Group

	vistorRouter := router.NewVistorRouter()
	group.Go(func() error {
		return vistorRouter.Run(fmt.Sprintf(":%v", global.Conf.App.Port))
	})
	log.Println("短链接服务启动")

	adminRouter := router.NewAdminRouter()
	group.Go(func() error {
		return adminRouter.Run(fmt.Sprintf(":%v", global.Conf.App.AdminPort))
	})
	log.Println("管理员HTTP服务启动")

	group.Go(func() error {
		storage.StartConsume(ctx, pool, global.StreamName, global.GroupName, global.ConsumerName)
		return nil
	})

	group.Go(func() error {
		DailyStatsWork(pool)
		return nil
	})

	group.Go(func() error {
		StatsWork(pool)
		return nil
	})

	err = group.Wait()
	log.Println("Group Failed:", err)
}

// DailyStatsWork 更新短链接每日数据表
func DailyStatsWork(pool *pgxpool.Pool) {
	ticker := time.NewTicker(global.UpdateTime)
	defer ticker.Stop()

	for range ticker.C {
		storage.SyncDailyStatsToPostgres(ctx,pool)
	}
}


// StatsWork 根据短链接每日数据表更新全表数据
func StatsWork(pool *pgxpool.Pool) {
	ticker := time.NewTicker(global.UpdateTime2)
	defer ticker.Stop()

	for range ticker.C {
		storage.SyncStatsToPostgres(ctx)
		storage.SyncTotalStatsToPostgres(ctx, pool)
	}
}
