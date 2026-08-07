package storage

import (
	"context"
	"log"
	"my_url_shortner/model"
	"my_url_shortner/utils"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type BatchFlusher struct {
	pgPool        *pgxpool.Pool
	stream        string
	group         string
	consumer      string
	batchSize     int
	flushInterval time.Duration
}

func ProduceMessage(streamName string, shortUrl, Ip string) {
	err := Rs.redisClient.XAdd(context.Background(), &redis.XAddArgs{
		Stream: streamName,
		MaxLen: 10000,
		Approx: true,
		ID:     "*",
		Values: map[string]interface{}{
			"short_url": shortUrl,
			"ip":        Ip,
		},
	}).Err()

	if err != nil {
		log.Println("Stream发送消息出错", err.Error())
	}
}

func StartConsume(ctx context.Context, pool *pgxpool.Pool, streamName, groupName, consumerName string) {
	err := Rs.redisClient.XGroupCreateMkStream(ctx, streamName, groupName, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		log.Println("消费者启动失败:", err)
		return
	}

	bf := &BatchFlusher{
		pgPool:        pool,
		stream:        streamName,
		group:         groupName,
		consumer:      consumerName,
		batchSize:     1000,
		flushInterval: 1 * time.Minute,
	}

	buffer := make([]model.AccessLog, 0, bf.batchSize)
	msgIDs := make([]string, 0, bf.batchSize)
	ticker := time.NewTicker(bf.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if len(buffer) > 0 {
				bf.flushToPostgres(context.Background(), buffer, msgIDs)
			}
			return
		case <-ticker.C:
			if len(buffer) > 0 {
				if err := bf.flushToPostgres(ctx, buffer, msgIDs); err == nil {
					buffer = buffer[:0]
					msgIDs = msgIDs[:0]
				} else {
					log.Println("更新日志出错!:", err.Error())
				}
			}
		default:
			// 从 Stream 拉取消息
			entries, err := Rs.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    bf.group,
				Consumer: bf.consumer,
				Streams:  []string{bf.stream, ">"},
				Count:    int64(bf.batchSize - len(buffer)),
				Block:    100 * time.Millisecond,
			}).Result()

			if err != nil || len(entries) == 0 {
				continue
			}

			for _, xmsg := range entries[0].Messages {
				url, _ := xmsg.Values["short_url"].(string)
				ip, _ := xmsg.Values["ip"].(string)

				if utils.EmptyString(url) || utils.EmptyString(ip) {
					_ = Rs.redisClient.XAck(ctx, bf.stream, bf.group, xmsg.ID).Err()
					continue
				}

				access_logs := model.AccessLog{
					ShortUrl: url,
					Ip:       ip,
				}
				buffer = append(buffer, access_logs)
				msgIDs = append(msgIDs, xmsg.ID)
			}

			// 达到批次大小，立即刷新
			if len(buffer) >= bf.batchSize {
				if err := bf.flushToPostgres(ctx, buffer, msgIDs); err == nil {
					buffer = buffer[:0]
					msgIDs = msgIDs[:0]
				} else {
					log.Println("更新日志出错!:", err.Error())
				}
			}
		}
	}

}

func (bf *BatchFlusher) flushToPostgres(ctx context.Context, buffer []model.AccessLog, msgsID []string) error {
	tx, err := bf.pgPool.Begin(ctx)
	if err != nil {
		log.Printf("开启 PG 事务失败: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	// 构建 CopyFromSource 迭代器数据
	rows := make([][]interface{}, len(buffer))
	for i, l := range buffer {
		rows[i] = []interface{}{l.ShortUrl, l.Ip}
	}

	// 利用 PG COPY 协议（吞吐量远超普通 INSERT）
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"access_logs"},
		[]string{"short_url", "ip"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		log.Printf("CopyFrom 批量写入 DB 失败: %v", err)
		return err // 写入失败时不确认 Redis ACK，下次重启或重试会重新读取
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("提交 PG 事务失败: %v", err)
		return err
	}

	// 数据库写入成功后，批量确认 Redis 消息
	return Rs.redisClient.XAck(ctx, bf.stream, bf.group, msgsID...).Err()
}
