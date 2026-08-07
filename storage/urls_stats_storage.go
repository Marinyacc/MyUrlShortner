package storage

import (
	"context"
	"fmt"
	"log"
	"my_url_shortner/global"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Ip_Count_Incr 用户点击短链接后更新Redis
func Ip_Count_Incr(shortUrl string, Ip string) {
	ctx := context.Background()
	today := time.Now().Format("20060102")

	pvDailyKey := global.ShortUrl_CountPrefix + today + ":" + shortUrl
	uvDailyKey := global.ShortUrl_IP_CountPrefix + today + ":" + shortUrl
	setKey := global.ShortUrl_Active_Set + today //加入到当天活跃的短链接set

	pipe := Rs.redisClient.Pipeline()

	pipe.Incr(ctx, pvDailyKey)
	pipe.PFAdd(ctx, uvDailyKey, Ip)
	pipe.SAdd(ctx, setKey, shortUrl)

	pipe.Incr(ctx, global.PvTotalKey)
	pipe.PFAdd(ctx, global.UvTotalKey, Ip)

	//设置每日 key的过期时间
	pipe.Expire(ctx, pvDailyKey, 8*24*time.Hour)
	pipe.Expire(ctx, uvDailyKey, 8*24*time.Hour)
	pipe.Expire(ctx, setKey, 2*24*time.Hour)

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Println("Ip_Count_Incr出错!:", err.Error())
	}
}

// SyncDailyStatsToPostgres 异步将当天更新的短链接统计数据加载到数据库
func SyncDailyStatsToPostgres(ctx context.Context, pool *pgxpool.Pool) {
	todayRedis := time.Now().Format("20060102")
	todayDb := time.Now().Format("2006-01-02") // 统一使用带连字符的标准日期格式存入数据库
	setKey := global.ShortUrl_Active_Set + todayRedis

	shortUrls, err := Rs.redisClient.SMembers(ctx, setKey).Result()
	if err != nil {
		log.Println("SyncDailyStatsToPostgres出错!:", err.Error())
		return
	}
	if len(shortUrls) == 0 {
		return
	}

	pipe := Rs.redisClient.Pipeline()
	pvCmds := make([]*redis.StringCmd, len(shortUrls))
	uvCmds := make([]*redis.IntCmd, len(shortUrls))

	for i, shortUrl := range shortUrls {
		pvDailyKey := global.ShortUrl_CountPrefix + todayRedis + ":" + shortUrl
		uvDailyKey := global.ShortUrl_IP_CountPrefix + todayRedis + ":" + shortUrl

		pvCmds[i] = pipe.Get(ctx, pvDailyKey)
		uvCmds[i] = pipe.PFCount(ctx, uvDailyKey)
	}

	_, err = pipe.Exec(ctx)

	if err != nil && err != redis.Nil {
		log.Println("SyncDailyStatsToPostgres出错!:", err.Error())
		return
	}

	batch := &pgx.Batch{}
	query := `INSERT INTO daily_stats (short_url, date, pv, uv)
              VALUES ($1, $2, $3, $4)
              ON CONFLICT (short_url, date)
              DO UPDATE SET pv = EXCLUDED.pv, uv = EXCLUDED.uv`

	for i, shortUrl := range shortUrls {
		pv, _ := pvCmds[i].Int()
		uv, _ := uvCmds[i].Result()
		batch.Queue(query, shortUrl, todayDb, pv, int(uv))
	}

	br := pool.SendBatch(ctx, batch)
	defer br.Close()
	for range shortUrls {
		if _, err := br.Exec(); err != nil {
			log.Println("SyncDailyStatsToPostgres出错!:", err.Error())
		}
	}
}

// SyncStatsToPostgres 定时任务
func SyncStatsToPostgres(ctx context.Context) {
	today := time.Now().Format("20060102")
	setKey := global.ShortUrl_Active_Set + today

	shortUrls, err := Rs.redisClient.SMembers(ctx, setKey).Result()
	if err != nil {
		log.Println("SyncStatsToPostgres出错!:", err.Error())
		return
	}
	if len(shortUrls) == 0 {
		return
	}

	for _, shortUrl := range shortUrls {
		RefreshShortUrlStats(ctx, shortUrl)
	}
}

// RefreshShortUrlStats 更新单链接统计表
func RefreshShortUrlStats(ctx context.Context, shortUrl string) {
	now := time.Now()
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	last7Days := now.AddDate(0, 0, -6).Format("2006-01-02")
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	query := `
        INSERT INTO stats (
            short_url, today_count, yesterday_count, last_7_days_count, monthly_count, total_count,
            d_today_count, d_yesterday_count, d_last_7_days_count, d_monthly_count, d_total_count
        )
        SELECT 
            @short_url::varchar AS short_url,
            COALESCE(SUM(CASE WHEN date = @today::date THEN pv END), 0),
            COALESCE(SUM(CASE WHEN date = @yesterday::date THEN pv END), 0),
            COALESCE(SUM(CASE WHEN date >= @last7Days::date THEN pv END), 0),
            COALESCE(SUM(CASE WHEN date >= @monthStart::date THEN pv END), 0),
            COALESCE(SUM(pv), 0),
            COALESCE(SUM(CASE WHEN date = @today::date THEN uv END), 0),
            COALESCE(SUM(CASE WHEN date = @yesterday::date THEN uv END), 0),
            COALESCE(SUM(CASE WHEN date >= @last7Days::date THEN uv END), 0),
            COALESCE(SUM(CASE WHEN date >= @monthStart::date THEN uv END), 0),
            COALESCE(SUM(uv), 0)
        FROM daily_stats
        WHERE short_url = @short_url::varchar
        ON CONFLICT (short_url) 
        DO UPDATE SET
            today_count = EXCLUDED.today_count,
            yesterday_count = EXCLUDED.yesterday_count,
            last_7_days_count = EXCLUDED.last_7_days_count,
            monthly_count = EXCLUDED.monthly_count,
            total_count = EXCLUDED.total_count,
            d_today_count = EXCLUDED.d_today_count,
            d_yesterday_count = EXCLUDED.d_yesterday_count,
            d_last_7_days_count = EXCLUDED.d_last_7_days_count,
            d_monthly_count = EXCLUDED.d_monthly_count,
            d_total_count = EXCLUDED.d_total_count;
    `

	err := DbNamedExec(ctx, query, pgx.NamedArgs{
		"short_url":  shortUrl,
		"today":      today,
		"yesterday":  yesterday,
		"last7Days":  last7Days,
		"monthStart": monthStart,
	})

	if err != nil {
		log.Println("RefreshShortUrlStats出错!:", err.Error())
	}
}

// SyncTotalStatsToPostgres 定时任务
func SyncTotalStatsToPostgres(ctx context.Context, pool *pgxpool.Pool) {
	err := RefreshStatsSumInGo(ctx, pool)
	if err != nil {
		log.Println("SyncTotalStatsToPostgre出错:", err)
	}
}

// RefreshStatsSumInGo 更新全表统计
func RefreshStatsSumInGo(ctx context.Context, pool *pgxpool.Pool) error {
	now := time.Now()
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	last7Days := now.AddDate(0, 0, -6).Format("2006-01-02")
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	total_pv_, err := RedisGetString(ctx, global.PvTotalKey)
	if err != nil {
		log.Println("RefreshStatsSumInGo出错!:", err.Error())
		return err
	}
	total_uv, err := Rs.redisClient.PFCount(ctx, global.UvTotalKey).Result()
	if err != nil {
		log.Println("RefreshStatsSumInGo出错!:", err.Error())
		return err
	}
	total_pv, _ := strconv.ParseInt(total_pv_, 10, 64)
	statsMap := map[string]int{
		"total_pv":            int(total_pv),
		"total_uv":            int(total_uv),
		"today_count":         0,
		"yesterday_count":     0,
		"last_7_days_count":   0,
		"monthly_count":       0,
		"d_today_count":       0,
		"d_yesterday_count":   0,
		"d_last_7_days_count": 0,
		"d_monthly_count":     0,
	}

	minQueryDate := monthStart
	if last7Days < minQueryDate {
		minQueryDate = last7Days
	}
	query := `
		SELECT 
			date::date, 
			COALESCE(SUM(pv), 0), 
			COALESCE(SUM(uv), 0) 
		FROM daily_stats 
		WHERE date >= $1::date 
		GROUP BY date
	`
	rows, err := pool.Query(ctx, query, minQueryDate)
	if err != nil {
		return fmt.Errorf("RefreshStatsSumInGo出错!: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var datee time.Time
		var pv, uv int
		if err := rows.Scan(&datee, &pv, &uv); err != nil {
			return err
		}

		date := datee.Format("2006-01-02")

		// 累加 PV/UV
		if date == today {
			statsMap["today_count"] += pv
			statsMap["d_today_count"] += uv
		}
		if date == yesterday {
			statsMap["yesterday_count"] += pv
			statsMap["d_yesterday_count"] += uv
		}
		if date >= last7Days {
			statsMap["last_7_days_count"] += pv
			statsMap["d_last_7_days_count"] += uv
		}
		if date >= monthStart {
			statsMap["monthly_count"] += pv
			statsMap["d_monthly_count"] += uv
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	upsertSQL := `
		INSERT INTO stats_sum (stats_key, stats_value) 
		VALUES ($1, $2)
		ON CONFLICT (stats_key) 
		DO UPDATE SET stats_value = EXCLUDED.stats_value
	`

	// 3. 开启事务进行批量更新
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("RefreshStatsSumInGo出错!: %w", err)
	}
	defer tx.Rollback(ctx) // 若已 Commit，Rollback 会自动被安全忽略

	// 4. 构建 pgx.Batch
	batch := &pgx.Batch{}
	for key, val := range statsMap {
		batch.Queue(upsertSQL, key, val)
	}

	// 5. 发送批处理并校验结果
	br := tx.SendBatch(ctx, batch)
	for range statsMap {
		if _, err := br.Exec(); err != nil {
			br.Close()
			return fmt.Errorf("RefreshStatsSumInGo出错!: %w", err)
		}
	}
	if err := br.Close(); err != nil {
		return fmt.Errorf("RefreshStatsSumInGo出错!: %w", err)
	}
	return tx.Commit(ctx)
}
