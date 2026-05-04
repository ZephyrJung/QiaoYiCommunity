# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

QiaoYi Community - Golang backend API for a residential community interaction system.

### Features
- WeChat mini-program login + phone SMS login
- Second-hand trading marketplace
- Property issue feedback system
- Community forum (posts, comments, likes)

### Tech Stack
- Go 1.24+
- Gin (web framework)
- GORM + MySQL
- Redis (SMS codes, tokens, rate limiting)
- JWT authentication
- Snowflake IDs
- Swagger docs

## Build & Run

### 本地开发

```bash
# Install dependencies
go mod tidy

# Run the server
go run cmd/server/main.go

# Build binary
go build -o server cmd/server/main.go

# Run tests
go test ./...

# Generate swagger docs
swag init -g cmd/server/main.go
```

Server starts on port 8080 by default (configurable in config/config.yaml).

### Docker 运行

使用 Docker Compose 一键启动所有服务（包含 MySQL、Redis 和应用）：

```bash
# 启动所有服务（后台运行）
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

**服务访问地址：**
- API服务: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/swagger/index.html`
- MySQL: `localhost:3306` (用户名: root, 密码: root, 数据库: qiaoyi_community)
- Redis: `localhost:6379`

**环境变量配置（在 docker-compose.yml 中修改）：**

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| DB_HOST | 数据库主机 | mysql |
| DB_PORT | 数据库端口 | 3306 |
| DB_USER | 数据库用户名 | root |
| DB_PASSWORD | 数据库密码 | root |
| DB_NAME | 数据库名称 | qiaoyi_community |
| REDIS_HOST | Redis主机 | redis |
| REDIS_PORT | Redis端口 | 6379 |
| SERVER_PORT | 服务端口 | :8080 |
| SERVER_MODE | 运行模式 | debug |

## Architecture

Layered architecture:
- `model/entity/` - GORM models
- `repository/dao/` - Data access layer (interfaces + implementations)
- `service/` - Business logic layer
- `handler/api/` - HTTP handlers (Gin)
- `middleware/` - Auth, CORS, recovery, rate limiting
- `pkg/` - Reusable packages (JWT, Redis, SMS, WeChat, OSS, response)
- `utils/` - Helpers (snowflake IDs, pagination, validation)
- `config/` - Viper-based configuration
- `migrations/` - SQL schema
- `docs/` - Swagger documentation

## Project Progress

### Completed (2026-05-04)
- **Project skeleton**: go.mod, directory structure, dependencies installed
- **Config layer**: `config/config.go`, `config/config.yaml`
- **Base packages**: `pkg/response/response.go`, `utils/snowflake.go`, `utils/paginate.go`, `utils/validator.go`
- **Infrastructure**:
  - `pkg/jwt/jwt.go` - JWT token generation/validation (access + refresh)
  - `pkg/redis/redis.go` - Redis client wrapper + key helpers
  - `pkg/sms/sms.go` - SMS provider interface + mock implementation
  - `pkg/wechat/miniprogram.go` - WeChat mini-program jscode2session
  - `pkg/oss/oss.go` + `pkg/oss/local_storage.go` - Local file upload storage
- **Database models** (7 entities):
  - `model/entity/user.go` - User, UserProfile, UserRole, UserStatus
  - `model/entity/category.go` - Category, CategoryType
  - `model/entity/item.go` - Item, ItemImage, ItemStatus
  - `model/entity/feedback.go` - Feedback, FeedbackImage, FeedbackReply, FeedbackType, FeedbackStatus
  - `model/entity/post.go` - Post, PostImage, PostStatus
  - `model/entity/comment.go` - Comment, CommentStatus
  - `model/entity/like.go` - Like, LikeTargetType
- **Migration script**: `migrations/001_init.up.sql` - 12 tables + 10 seeded categories
- **DAO layer** (7 DAOs with interfaces):
  - `repository/dao/user_dao.go` - Create, GetByID, GetByUnionID, GetByOpenID, GetByPhone, Update, GetProfile, SaveProfile
  - `repository/dao/category_dao.go` - ListByType, GetByID
  - `repository/dao/item_dao.go` - Create, GetByID, Update, Delete (soft), List with filters
  - `repository/dao/feedback_dao.go` - Create, GetByID, Update, Delete (soft), List, CreateReply
  - `repository/dao/post_dao.go` - Create, GetByID, Update, Delete (soft), List, IncrementViewCount
  - `repository/dao/comment_dao.go` - Create, GetByID, Delete (soft), ListByPostID
  - `repository/dao/like_dao.go` - Create, Delete, Exists
- **Middleware**:
  - `middleware/jwt_auth.go` - JWT authentication + RequireRole
  - `middleware/cors.go` - CORS headers for cross-origin requests
  - `middleware/recovery.go` - Panic recovery with structured logging
  - `middleware/rate_limit.go` - Redis-based sliding window rate limiter
- **Entry Point & Routing**:
  - `cmd/server/main.go` - Config loading, DB/Redis init, dependency injection, Gin setup, graceful shutdown
  - `router/router.go` - Route grouping: public routes, auth-protected routes
- **Service Layer** (8 services):
  - `service/auth_service.go` - WeChat login, phone SMS login, token refresh, logout
  - `service/user_service.go` - Profile get/update
  - `service/item_service.go` - CRUD + search + ownership checks + mark sold
  - `service/feedback_service.go` - CRUD + reply ACL + status transitions
  - `service/post_service.go` - CRUD + view counting
  - `service/comment_service.go` - Create + tree listing
  - `service/like_service.go` - Atomic toggle + counter sync
  - `service/upload_service.go` - File validation (jpeg/png, max 5MB, max 9 files) + storage delegation
- **Handler Layer** (8 handlers):
  - `handler/api/auth_handler.go` - Login endpoints
  - `handler/api/user_handler.go` - Profile endpoints
  - `handler/api/item_handler.go` - Item endpoints
  - `handler/api/feedback_handler.go` - Feedback endpoints
  - `handler/api/post_handler.go` - Post endpoints
  - `handler/api/comment_handler.go` - Comment endpoints
  - `handler/api/like_handler.go` - Like toggle/status endpoints
  - `handler/api/upload_handler.go` - Multipart image upload
  - `handler/api/category_handler.go` - Category listing
- **Build & Docs**:
  - `Makefile` - run, build, test, migrate, swagger targets
  - `docs/` - Swagger docs generated via `swag init`

### Verified Status
- ✅ `go build ./...` - Compilation successful
- ✅ `go test ./...` - Tests pass
- ✅ `swag init -g cmd/server/main.go` - Swagger docs generated

### Remaining Tasks (Optional Enhancements)
- [ ] Add unit tests for critical services
- [ ] Add Swagger annotations to handlers for better API documentation
- [ ] Implement admin/manager role endpoints
- [ ] Add production configuration (env-specific configs)
- [ ] Add logging middleware
- [ ] Implement proper error handling with custom error types

## API Endpoints

### Public Routes (No Auth Required)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/wechat/login` | WeChat mini-program login |
| POST | `/api/v1/auth/phone/send-code` | Send SMS verification code |
| POST | `/api/v1/auth/phone/login` | Phone SMS login |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| GET | `/api/v1/categories` | List categories |
| GET | `/api/v1/items` | List items (public) |
| GET | `/api/v1/items/:id` | Get item detail |
| GET | `/api/v1/posts` | List posts (public) |
| GET | `/api/v1/posts/:id` | Get post detail |
| GET | `/api/v1/comments/post/:post_id` | List comments by post |

### Protected Routes (JWT Required)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth/logout` | Logout user |
| GET | `/api/v1/user/profile` | Get user profile |
| PUT | `/api/v1/user/profile` | Update user profile |
| POST | `/api/v1/items` | Create item |
| PUT | `/api/v1/items/:id` | Update item |
| DELETE | `/api/v1/items/:id` | Delete item |
| PATCH | `/api/v1/items/:id/sold` | Mark item as sold |
| POST | `/api/v1/feedbacks` | Create feedback |
| GET | `/api/v1/feedbacks` | List feedbacks |
| GET | `/api/v1/feedbacks/:id` | Get feedback detail |
| PUT | `/api/v1/feedbacks/:id` | Update feedback |
| DELETE | `/api/v1/feedbacks/:id` | Delete feedback |
| POST | `/api/v1/feedbacks/:id/replies` | Reply to feedback |
| POST | `/api/v1/posts` | Create post |
| PUT | `/api/v1/posts/:id` | Update post |
| DELETE | `/api/v1/posts/:id` | Delete post |
| POST | `/api/v1/comments` | Create comment |
| DELETE | `/api/v1/comments/:id` | Delete comment |
| POST | `/api/v1/likes/toggle` | Toggle like status |
| GET | `/api/v1/likes/status` | Check like status |
| POST | `/api/v1/upload/images` | Upload images |