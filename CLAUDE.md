# Jump Jump - 短链接系统

## 项目概述

Jump Jump 是一个用 Go 语言开发的功能完善的短链接（URL Shortener）系统，支持短链接的创建、管理、访问统计和报表生成。

## 技术栈

- **语言**: Go 1.13
- **Web 框架**: Gin (v1.6.3)
- **数据存储**: Redis（唯一存储，无关系型数据库）
- **配置管理**: [reborn](https://github.com/jwma/reborn)（基于 Redis 的动态配置库）
- **认证**: JWT (dgrijalva/jwt-go) + scrypt 密码哈希
- **API 文档**: Swagger (swaggo/swag + gin-swagger)
- **前端管理后台**: Vue.js + Element UI（预编译的静态资源，源码不在本仓库）
- **容器化**: Docker + docker-compose
- **CI/CD**: GitHub Actions

## 项目架构

```
cmd/                        # 可执行程序入口
├── apiserver/              # API 服务 + 管理后台（默认 :8080）
├── landingserver/          # 短链接跳转服务（默认 :8081）
├── createuser/             # CLI 工具：创建用户
├── reportgenerator/        # 后台定时任务：日报生成器（每 30s）
└── migraterequesthistory/  # 数据迁移工具（v1.2.0 / v1.3.0）

internal/app/
├── cmd/server/             # 服务启动逻辑（初始化 DB、Config、Router）
├── config/                 # 系统配置管理（域名、ID 长度、404 处理）
├── db/                     # Redis 连接管理
├── handlers/               # HTTP 请求处理器
│   ├── middleware.go       # JWT 认证中间件 + AllowedHosts 中间件
│   ├── user.go             # 登录/登出/修改密码
│   ├── shortlink.go        # 短链接 CRUD + 数据查询
│   ├── config.go           # 系统配置 API
│   └── landing.go          # 短链接跳转（307 重定向）
├── models/                 # 数据模型（User, ShortLink, RequestHistory, DailyReport）
├── repository/             # 数据访问层（Redis 操作封装）
├── report/                 # 日报生成逻辑
│   ├── service.go          # PV/UV/OS 统计计算
│   └── generator.go        # 定时任务调度器
├── routers/                # 路由注册
└── utils/                  # 工具函数（密钥、字符串、Redis Key 生成）

web/admin/                  # 管理后台前端（预编译的 Vue.js SPA）
deployments/                # docker-compose 部署配置
build/package/              # Dockerfile（多阶段构建）
docs/                       # Swagger 文档（自动生成）
```

## 双服务架构

系统由两个独立服务组成：

1. **apiserver** (`cmd/apiserver`): 提供 REST API 和管理后台 UI
   - 端口: 由 `J2_API_ADDR` 环境变量控制
   - 路由前缀: `/v1/`
   - 同时托管管理后台静态资源 (`web/admin/`)

2. **landingserver** (`cmd/landingserver`): 处理短链接跳转
   - 端口: 由 `J2_LANDING_ADDR` 环境变量控制
   - 访问 `/:id` → 307 重定向到目标 URL
   - 异步记录访问历史和活跃链接

## API 路由

### 用户相关
- `POST /v1/user/login` — 登录
- `GET /v1/user/info` — 获取用户信息（需 JWT）
- `POST /v1/user/logout` — 登出（需 JWT）
- `POST /v1/user/change-password` — 修改密码（需 JWT）

### 系统配置
- `GET /v1/config` — 获取系统配置（需 JWT）
- `PATCH /v1/config/landing-hosts` — 更新域名（管理员）
- `PATCH /v1/config/id-length` — 更新 ID 长度设置（管理员）
- `PATCH /v1/config/short-link-404-handling` — 更新 404 处理（管理员）

### 短链接
- `GET /v1/short-link/` — 短链接列表（分页）
- `GET /v1/short-link/:id` — 获取短链接详情
- `POST /v1/short-link/` — 创建短链接
- `PATCH /v1/short-link/:id` — 更新短链接
- `DELETE /v1/short-link/:id` — 删除短链接
- `GET /v1/short-link/:id/data` — 获取访问数据（按日期范围）

## 数据模型

### User
- 字段: Username, Role(1=普通用户, 2=管理员), Password(scrypt), Salt, CreateTime
- 存储: Redis Hash (`users`)

### ShortLink
- 字段: Id, Url, Description, IsEnable, CreatedBy, CreateTime, UpdateTime
- 存储: Redis String (`link:{id}`)，索引用 Sorted Set (`links`, `links:{username}`)

### RequestHistory
- 字段: Id, Url, IP, UA, Time
- 存储: Redis Sorted Set (`rh:{linkId}`)，score 为时间戳

### DailyReport
- 字段: PV, UV, OS(按操作系统统计)
- 存储: Redis Hash (`dr:{linkId}`)，field 为日期字符串

## 权限模型

- **管理员 (Role=2)**: 可查看所有短链接、指定自定义 ID、修改系统配置
- **普通用户 (Role=1)**: 只能查看/管理自己创建的短链接，ID 自动生成
- 所有 API（除登录外）需要 JWT 认证

## 环境变量

| 变量 | 说明 | 示例 |
|------|------|------|
| `J2_API_ADDR` | API 服务监听地址 | `0.0.0.0:8080` |
| `J2_LANDING_ADDR` | 跳转服务监听地址 | `0.0.0.0:8081` |
| `REDIS_HOST` | Redis 地址 | `db:6379` |
| `REDIS_DB` | Redis 数据库编号 | `0` |
| `REDIS_PASSWORD` | Redis 密码 | |
| `SECRET_KEY` | JWT 签名密钥 | |
| `GIN_MODE` | Gin 运行模式 | `debug` / `release` |
| `ALLOWED_HOSTS` | 允许的 Host 列表（release 模式必填） | `localhost,127.0.0.1` |
| `API_DOC_USERNAME` | Swagger 文档用户名 | |
| `API_DOC_PASSWORD` | Swagger 文档密码 | |
| `LANDING_HOME` | Landing 首页跳转地址 | |

## 本地开发

```bash
# 启动依赖服务（Redis）
docker-compose -f deployments/docker-compose.yaml -p jumpjump up -d

# 创建用户
docker-compose -f deployments/docker-compose.yaml -p jumpjump exec apiserver ./createuser -username=admin -password=12345 -role=2

# 访问管理后台
# http://localhost:8080
```

## 构建与部署

```bash
# 构建 Docker 镜像
make dockerimage

# 生成 Swagger 文档
make docs
```

## Redis Key 规范

| Key 模式 | 类型 | 说明 |
|----------|------|------|
| `users` | Hash | 用户数据，field 为 username |
| `link:{id}` | String (JSON) | 短链接详情 |
| `links` | Sorted Set | 全局短链接索引 |
| `links:{username}` | Sorted Set | 用户短链接索引 |
| `rh:{linkId}` | Sorted Set | 访问历史，score 为时间戳 |
| `activelinks` | Sorted Set | 活跃链接，score 为时间戳 |
| `dr:{linkId}` | Hash | 日报数据，field 为日期 |
| `j2config` | - | 系统配置（reborn 管理） |

## Agent skills

### Issue tracker

Issues live in GitHub (jwma/jump-jump) via the `gh` CLI. See `docs/agents/issue-tracker.md`.

### Triage labels

Five canonical triage roles, each mapped to its own label (needs-triage, needs-info, ready-for-agent, ready-for-human, wontfix). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context repo — one `CONTEXT.md` + `docs/adr/` at the repo root. See `docs/agents/domain.md`.
