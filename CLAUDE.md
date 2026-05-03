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
- Go 1.22+
- Gin (web framework)
- GORM + MySQL
- Redis (SMS codes, tokens, rate limiting)
- JWT authentication
- Snowflake IDs
- Swagger docs

## Build & Run

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

## Project Progress

### Completed (2026-05-03)
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
  - `repository/dao/user_dao.go` - Create, GetByID, GetByUnionID, GetByPhone, Update
  - `repository/dao/category_dao.go` - ListByType, GetByID
  - `repository/dao/item_dao.go` - Create, GetByID, Update, Delete (soft), List with filters
  - `repository/dao/feedback_dao.go` - Create, GetByID, Update, Delete (soft), List, CreateReply
  - `repository/dao/post_dao.go` - Create, GetByID, Update, Delete (soft), List, IncrementViewCount
  - `repository/dao/comment_dao.go` - Create, GetByID, Delete (soft), ListByPostID
  - `repository/dao/like_dao.go` - Create, Delete, Exists
- **Auth middleware**: `middleware/jwt_auth.go` + `RequireRole` middleware

### Missing / Remaining Files

#### Middleware (2 of 3 done conceptually, only jwt_auth.go exists)
- [ ] `middleware/cors.go` - CORS headers for cross-origin requests
- [ ] `middleware/recovery.go` - Panic recovery with structured logging
- [ ] `middleware/rate_limit.go` - Redis-based sliding window rate limiter

#### Entry Point & Routing
- [ ] `cmd/server/main.go` - Config loading, DB/Redis init, dependency injection, Gin setup, graceful shutdown
- [ ] `router/router.go` - Route grouping: public routes, auth-protected routes, manager-protected routes

#### Service Layer (8 services)
- [ ] `service/auth_service.go` - WeChat login, phone SMS login, token refresh, logout
- [ ] `service/user_service.go` - Profile get/update
- [ ] `service/item_service.go` - CRUD + search + ownership checks + mark sold
- [ ] `service/feedback_service.go` - CRUD + reply ACL + status transitions
- [ ] `service/post_service.go` - CRUD + view counting
- [ ] `service/comment_service.go` - Create + tree listing
- [ ] `service/like_service.go` - Atomic toggle + counter sync in transaction
- [ ] `service/upload_service.go` - File validation (jpeg/png, max 5MB, max 9 files) + storage delegation

#### Handler Layer (8 handlers + DTOs)
- [ ] `handler/api/auth_handler.go` - Login endpoints + Swagger annotations
- [ ] `handler/api/user_handler.go` - Profile endpoints
- [ ] `handler/api/item_handler.go` - Item endpoints
- [ ] `handler/api/feedback_handler.go` - Feedback endpoints
- [ ] `handler/api/post_handler.go` - Post endpoints
- [ ] `handler/api/comment_handler.go` - Comment endpoints
- [ ] `handler/api/like_handler.go` - Like toggle/status endpoints
- [ ] `handler/api/upload_handler.go` - Multipart image upload

#### Build & Docs
- [ ] `Makefile` - run, build, test, migrate, swagger targets
- [ ] `docs/` - Swagger docs generated via `swag init`
- [ ] Swagger annotations on all handlers

#### Potential Dependency Gap
- `pkg/jwt/jwt.go` imports `github.com/google/uuid` which may need `go get`

## Next Steps to Resume
1. Finish middleware: cors.go, recovery.go, rate_limit.go
2. Create cmd/server/main.go with full DI wiring
3. Create router/router.go with route groups
4. Implement service/auth_service.go + handler/api/auth_handler.go
5. Implement remaining services and handlers in order: upload/category, item, feedback, post, comment, like
6. Add Swagger annotations to all handlers, run `swag init`
7. Create Makefile
8. Update go.mod if uuid dependency is missing
