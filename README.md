# Jump Jump

> 开箱即用的短链接管理平台，支持多租户

## 特性

- 多租户架构，支持团队协作与成员邀请
- 自定义域名绑定（每个租户可绑定多个域名）
- 短链接 CRUD + 启用/禁用
- 访问统计分析（PV/UV/OS 分布/IP 分布/访问趋势）
- 成员管理与角色权限（Admin / Member）
- 超级管理员后台（用户管理、租户管理）
- 二维码生成

## 技术栈

### 后端

- Go 1.24 + Gin
- PostgreSQL 17（主存储）
- Redis 7（缓存 + 访问记录缓冲）

### 前端（web-ui）

- Vue 3 + TypeScript
- Tailwind CSS 4 + Radix Vue
- Pinia + Vue Router
- ECharts（数据可视化）
- Vite 8

## 快速开始

### 前置条件

- Docker & Docker Compose

### Docker Compose（推荐）

```bash
# 克隆项目
git clone https://github.com/jwma/jump-jump.git
cd jump-jump

# 构建并启动所有服务（PostgreSQL、Redis、API Server、Landing Server）
docker compose -f docker-compose.dev.yaml up -d --build

# 查看服务状态
docker compose -f docker-compose.dev.yaml ps
```

### 创建超级管理员

服务启动后，运行以下命令创建超级管理员：

```bash
docker compose -f docker-compose.dev.yaml exec apiserver ./createsuperuser \
  -username=<用户名> \
  -password=<密码>
```

该命令需要通过环境变量 `DATABASE_URL` 连接数据库（容器内已自动配置）。

### 访问

- **管理后台**：http://localhost:8080
- **API 文档（Swagger）**：http://localhost:8080/swagger/index.html（默认账号：`apidoc` / `showmethedoc`）
- **API Server**：http://localhost:8080/v1/
- **短链接跳转**：http://localhost:8081/{id}

## 开发

### 前置条件

- Go 1.24+
- Node.js 20+
- PostgreSQL 17
- Redis 7

### 后端开发

```bash
# 设置环境变量（根据实际情况修改）
export DATABASE_URL="postgres://jumpjump:jumpjump@localhost:5432/jumpjump?sslmode=disable"
export REDIS_HOST="localhost:6379"
export REDIS_DB="0"
export SECRET_KEY="your-secret-key"
export J2_API_ADDR="0.0.0.0:8080"
export J2_LANDING_ADDR="0.0.0.0:8081"
export GIN_MODE="debug"
export ALLOWED_HOSTS="localhost,127.0.0.1"

# 启动 API Server
go run ./cmd/apiserver

# 启动 Landing Server（另开终端）
export J2_LANDING_ADDR="0.0.0.0:8081"
go run ./cmd/landingserver

# 创建超级管理员
go run ./cmd/createsuperuser -username=admin -password=yourpassword
```

> 开发环境可以使用 `docker compose -f docker-compose.dev.yaml up postgres redis -d` 单独启动 PostgreSQL 和 Redis，再用上述命令本地运行后端服务。

### 前端开发

```bash
cd web-ui
npm install
npm run dev
```

Vite 开发服务器已配置代理，`/v1` 请求会自动转发到 `http://localhost:8080`。

### 项目结构

```
.
├── cmd/
│   ├── apiserver/          # API Server 入口
│   ├── landingserver/      # Landing Server 入口（短链接跳转）
│   └── createsuperuser/    # 创建超级管理员命令
├── internal/app/
│   ├── cmd/server/         # 服务启动逻辑
│   ├── config/             # 租户配置管理
│   ├── db/                 # 数据库初始化与迁移
│   ├── handlers/           # HTTP 请求处理器
│   ├── models/             # 数据模型
│   ├── repository/         # 数据访问层
│   ├── routers/            # 路由定义
│   ├── utils/              # 工具函数（JWT 等）
│   └── workers/            # 后台任务（访问记录刷写）
├── web-ui/                 # 前端 Vue 3 应用
│   └── src/
│       ├── api/            # API 请求封装
│       ├── components/     # Vue 组件
│       ├── composables/    # 组合式函数
│       ├── i18n/           # 国际化
│       ├── layouts/        # 页面布局
│       ├── pages/          # 页面组件
│       ├── router/         # 路由配置
│       ├── stores/         # Pinia 状态管理
│       ├── types/          # TypeScript 类型
│       └── utils/          # 工具函数
├── docs/                   # Swagger 文档
├── build/package/          # Docker 构建文件
├── deployments/            # 部署配置（生产 docker-compose）
├── docker-compose.dev.yaml # 开发环境 Docker Compose
└── Makefile                # 构建命令
```

## 环境变量

| 变量 | 说明 | 必填 | 默认值 |
|------|------|------|--------|
| `DATABASE_URL` | PostgreSQL 连接串 | 是 | — |
| `REDIS_HOST` | Redis 地址 | 是 | — |
| `REDIS_DB` | Redis 数据库编号 | 是 | — |
| `REDIS_PASSWORD` | Redis 密码 | 否 | 空 |
| `SECRET_KEY` | JWT 签名密钥 | 是 | — |
| `J2_API_ADDR` | API Server 监听地址 | 是 | — |
| `J2_LANDING_ADDR` | Landing Server 监听地址 | 是 | — |
| `GIN_MODE` | Gin 运行模式（`debug` / `release`） | 否 | `debug` |
| `ALLOWED_HOSTS` | 允许的域名（逗号分隔，`release` 模式下必填） | 否 | — |
| `API_DOC_USERNAME` | Swagger 文档认证用户名 | 否 | `apidoc` |
| `API_DOC_PASSWORD` | Swagger 文档认证密码 | 否 | `showmethedoc` |
| `API_DOC_HOST` | Swagger 文档显示的主机地址 | 否 | 自动检测 |
| `PG_PASSWORD` | Docker 环境下 PostgreSQL 密码 | 否 | `jumpjump` |

## 部署

### 构建镜像

```bash
make dockerimage
```

### 使用 Docker Compose 部署

生产环境配置位于 `deployments/docker-compose.yaml`，使用 `release` 模式：

```bash
cd deployments

# 设置必要的密码
export PG_PASSWORD="your-secure-password"

# 启动服务
docker compose up -d
```

> 生产环境必须设置 `SECRET_KEY` 和 `ALLOWED_HOSTS`，否则 API Server 将拒绝启动。

### 创建超级管理员

```bash
docker compose exec apiserver ./createsuperuser -username=admin -password=<password>
```

## 常用命令

```bash
make dockerimage   # 构建 Docker 镜像
make docs          # 生成 Swagger 文档
make dev-ui        # 启动前端开发服务器
make build-ui      # 构建前端产物
make build         # 编译所有 Go 二进制
```

## License

[MIT](LICENSE)
