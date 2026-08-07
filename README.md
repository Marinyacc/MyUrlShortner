# MyUrlShortner V1

一个基于 Go (Gin) 构建的高性能短链接服务，支持 Web 管理后台和 RESTful API，使用 Redis 作为缓存层、PostgreSQL 作为持久化存储，并通过 Docker Compose 实现一键部署。

## ✨ 功能特性

- **短链接生成与跳转** — 输入长链接，自动生成唯一短链接；访问短链接时 301 重定向至目标地址
- **Web 管理后台** — 登录验证码保护，提供仪表盘、短链接管理、数据统计、访问日志等页面
- **RESTful API** — 完整的用户管理和短链接 CRUD 接口，便于第三方集成
- **PV / UV 实时统计** — 基于 Redis 计数器（PV）和 HyperLogLog（UV）实现高性能去重统计
(注意：UV只有全局统计是精确统计，其他的只是简单把统计数据相加，得到的结果不精确)
- **缓存穿透保护** — Redis 缓存未命中时使用分布式锁，防止大量请求击穿至数据库
- **Redis Streams 异步日志** — 访问日志通过 Redis Streams 异步写入 PostgreSQL，削峰解耦
- **定时统计同步** — 每分钟将 Redis 中的统计数据同步至 PostgreSQL，保证数据持久化
- **Docker 一键部署** — 提供完整的 Docker Compose 编排（server + PostgreSQL + Redis）

## 🛠️ 技术栈

| 类别 | 技术 |
|------|------|
| 编程语言 | Go 1.25 |
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) |
| 数据库 | [PostgreSQL](https://www.postgresql.org/) (pgx/v5) |
| 缓存 | [Redis](https://redis.io/) (go-redis/v9) |
| 认证 | JWT (golang-jwt/v5) |
| 验证码 | [dchest/captcha](https://github.com/dchest/captcha) |
| 密码加密 | bcrypt (golang.org/x/crypto) |
| 容器化 | Docker + Docker Compose |

## 📁 项目结构

```
.
├── main.go                  # 应用入口，启动 HTTP 服务和后台定时任务
├── config.ini               # 本地开发配置文件
├── config_docker.ini        # Docker 环境配置文件
├── .env                     # JWT Secret 等环境变量
├── structure.sql            # 数据库表结构初始化脚本
├── global/                  # 全局常量与变量定义
│   └── global.go
├── handler/                 # HTTP 请求处理器
│   ├── admin_handler.go     # 管理后台页面处理（登录、仪表盘等）
│   ├── urls_handler.go      # 短链接跳转处理
│   └── api_handler.go       # RESTful API 处理
├── middleware/               # Gin 中间件
│   ├── jwt.go               # JWT 认证中间件
│   └── log.go               # 日志中间件
├── model/                   # 数据模型与配置结构体
│   ├── config.go
│   ├── user.go
│   ├── short_url.go
│   ├── access_log.go
│   ├── stats.go
│   └── result.go
├── router/                  # 路由定义
│   └── router.go            # 访客路由（:9091）与管理员路由（:9092）
├── service/                 # 业务逻辑层
│   ├── user_service.go
│   └── urls_service.go
├── storage/                 # 数据访问层（PostgreSQL + Redis）
│   ├── postgre_init.go      # PostgreSQL 连接池初始化
│   ├── redis_init.go        # Redis 客户端初始化
│   ├── redis_stream.go      # Redis Streams 生产者 / 消费者
│   ├── redis_lock.go        # Redis 分布式锁
│   ├── captcha_storage.go   # 验证码 Redis 存储
│   ├── users_storage.go     # 用户数据操作
│   ├── urls_storage.go      # 短链接数据操作
│   ├── urls_stats_storage.go # 统计数据操作
│   ├── access_logs_storage.go # 访问日志操作
│   └── total_storage.go     # 汇总统计操作
├── templates/               # Go 模板（管理后台页面）
│   ├── login.gohtml
│   ├── dashboard.gohtml
│   ├── urls.gohtml
│   ├── urls_stats.gohtml
│   ├── access_logs.gohtml
│   └── error.gohtml
├── utils/                   # 工具函数
│   ├── config.go            # INI 配置文件解析
│   ├── rand.go              # 随机字符串生成
│   └── utils.go             # 通用工具函数
└── docker/                  # Docker 部署文件
    ├── Dockerfile
    ├── docker-compose.yaml
    ├── vars.env             # Docker 环境变量模板
    ├── start.sh             # 一键启动脚本
    └── clear.sh             # 清理脚本
```

## 🚀 快速开始

### 前置条件

- Go 1.25+
- PostgreSQL 16+
- Redis 7+
- Docker & Docker Compose（可选，容器化部署方式所需）

### Docker 部署

运行./docker目录下的start.sh

## 📡 API 接口

### 短链接跳转

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/:url` | 根据短链接跳转至目标长链接 |

### 管理后台（Web 页面）

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/login` | 登录页面 | 否 |
| POST | `/login` | 提交登录 | 否 |
| GET | `/captcha/:id` | 获取验证码图片 | 否 |
| GET | `/captcha` | 获取验证码 ID | 否 |
| POST | `/admin/logout` | 登出 | 是 |
| GET | `/admin/dashboard` | 仪表盘页面 | 是 |
| GET | `/admin/urls` | 短链接列表 | 是 |
| GET | `/admin/stats` | 数据统计列表 | 是 |
| GET | `/admin/access_logs` | 访问日志列表 | 是 |
| POST | `/admin/urls/generate` | 新建短链接 | 是 |
| POST | `/admin/urls/delete` | 删除短链接 | 是 |

### RESTful API

**用户管理**

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/account` | 获取所有用户信息 |
| POST | `/api/account` | 创建新用户 |
| POST | `/api/account/update` | 更新用户密码 |
| POST | `/api/account/delete` | 删除用户 |

**短链接管理**

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/urls` | 获取所有短链接信息 |
| GET | `/api/urls/:url` | 获取指定短链接信息 |
| POST | `/api/url` | 生成短链接 |
| GET | `/api/urls/stats` | 获取所有短链接统计 |
| GET | `/api/urls/stats/:url` | 获取指定短链接统计 |


## 🏗️ 系统架构

```
┌──────────┐    ┌──────────┐    ┌───────────┐
│  访客端   │    │  管理后台  │    │  API 客户端 │
│ :9091    │    │ :9092    │    │ :9092     │
└────┬─────┘    └────┬─────┘    └─────┬─────┘
     │               │               │
     ▼               ▼               ▼
┌──────────────────────────────────────────┐
│               Gin HTTP Server            │
│  ┌─────────┐  ┌──────────┐  ┌────────┐  │
│  │ Visitor  │  │  Admin   │  │  API   │  │
│  │ Router   │  │  Router   │  │ Router │  │
│  └────┬─────┘  └────┬─────┘  └───┬────┘  │
│       │             │            │       │
│  ┌────▼─────────────▼────────────▼────┐  │
│  │          Service / Handler         │  │
│  └────┬──────────────────┬───────────┘  │
└───────┼──────────────────┼──────────────┘
        │                  │
   ┌────▼────┐       ┌─────▼──────┐
   │  Redis   │       │ PostgreSQL │
   │ 缓存/锁   │       │ 持久化存储  │
   │ Streams  │       │            │
   └──────────┘       └────────────┘
```

### 短链接跳转流程

```
用户访问短链接
      │
      ▼
┌─────────────┐    未命中    ┌──────────────┐
│  Redis 缓存  │ ──────────► │ 分布式锁获取   │
│  查询长链接  │             │              │
└──────┬──────┘             └──────┬───────┘
       │ 命中                      │
       ▼                           ▼
   返回重定向              ┌──────────────┐
                          │ PostgreSQL   │
                          │ 查询长链接    │
                          └──────┬───────┘
                                 │
                                 ▼
                          ┌──────────────┐
                          │ 回写 Redis   │
                          │ 释放分布式锁  │
                          └──────────────┘
```

### 数据统计流程

```
短链接访问 → Redis 计数器 +1 (PV)
           → HyperLogLog 记录 IP (UV)
           → Redis Streams 写入访问日志
                       │
              ┌────────▼─────────┐
              │ Consumer 消费     │
              │ 批量写入 PostgreSQL│
              └──────────────────┘
                       │
              ┌────────▼─────────┐
              │ 定时任务（每分钟） │
              │ daily_stats →    │
              │ stats → stats_sum│
              └──────────────────┘
```

## 📊 数据库设计

| 表名 | 说明 |
|------|------|
| `users` | 管理员用户表，支持 bcrypt 密码加密 |
| `short_urls` | 短链接映射表（长链接 ↔ 短链接） |
| `access_logs` | 访问日志表，记录每次跳转的 IP 和时间 |
| `daily_stats` | 每日 PV/UV 统计表（按短链接 + 日期） |
| `stats` | 短链接多维统计表（今日/昨日/7天/月/累计） |
| `stats_sum` | 全局汇总统计（总 PV、总 UV 等） |

## 🔑 Redis 数据结构

| Key 模式 | 类型 | 说明 |
|----------|------|------|
| `Url_Count:20060102:shortUrl` | String (int) | 某短链接某日 PV |
| `Url_d_Count:20060102:shortUrl` | HyperLogLog | 某短链接某日 UV |
| `active_set:20060102` | Set | 当日活跃短链接集合 |
| `Url:short_url` | String | 短链接 → 长链接缓存 |
| `Url_Lock:short_url` | String | 短链接分布式锁 |
| `pvTotalKey` | String (int) | 全局 PV 计数 |
| `uvTotalKey` | HyperLogLog | 全局 UV 计数 |

## 📝 默认账号

首次初始化数据库后，系统会自动创建一个管理员账号：

- **账号**: `root`
- **密码**: `123456`

## 📄 License

MIT
