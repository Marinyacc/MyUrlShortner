package storage

import (
	"context"
	"fmt"
	"log"
	"my_url_shortner/global"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// postgreSQL 连接池
var dbPool *pgxpool.Pool

// InitPostgreService 初始化数据库服务
func InitPostgreService() (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", global.Conf.DataBase.User, global.Conf.DataBase.Password, global.Conf.DataBase.Host, global.Conf.DataBase.Port, global.Conf.DataBase.DbName)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConns = int32(global.Conf.DataBase.MaxOpenConns)                             //设置最大连接数
	config.MaxConnIdleTime = time.Duration(global.Conf.DataBase.MaxIdleTime) * time.Minute // 设置最大空闲时间

	log.Println("max_conn_life_time:", config.MaxConnLifetime)
	log.Println("max_open_conns:", config.MaxConns)
	log.Println("max_conn_idle_time:", config.MaxConnIdleTime)
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	// 增加 5 秒初始连接超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 增加 Ping 探活校验，确保服务启动时数据库可用
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database failed: %w", err)
	}

	if err != nil {
		return nil, err
	}
	dbPool = pool
	return dbPool, nil
}

// DbNamedExec 执行带有命名参数（@）的sql语句
func DbNamedExec(ctx context.Context, query string, args pgx.NamedArgs) error {
	_, err := dbPool.Exec(ctx, query, args)
	return err
}

// DbGet 获取单条记录
func DbGet(ctx context.Context, query string, dest any, args ...any) error {
	err := pgxscan.Get(ctx, dbPool, dest, query, args...)
	//是否忽略找不到数据的数据库错误
	// if err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return nil
	// 	}
	// }
	return err
}

// DbSelect 获取多条记录
func DbSelect(ctx context.Context, query string, dest any, args ...any) error {
	return pgxscan.Select(ctx, dbPool, dest, query, args...)
}

// DbClose 关闭连接池
func DbClose() {
	dbPool.Close()
}

type Query struct {
	Query string
	Args  []any
}

// DbExecTx 使用 pgxpool 执行一组事务 SQL
func DbExecTx(ctx context.Context, queries ...Query) error {
	// BeginFunc 会自动管理事务：
	// 1. 成功执行闭包内代码时自动 Commit
	// 2. 返回错误或发生 panic 时自动 Rollback
	err := pgx.BeginFunc(ctx, dbPool, func(tx pgx.Tx) error {
		for _, q := range queries {
			_, err := tx.Exec(ctx, q.Query, q.Args...)
			if err != nil {
				// 返回错误会触发自动回滚
				return err
			}
		}
		return nil
	})
	return err
}
