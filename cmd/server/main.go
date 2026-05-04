// @title QiaoYi Community API
// @version 1.0
// @description 桥驿社区互动系统后端 API（微信登录、二手交易、物业反馈、业主交流）
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {token}" to authenticate

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZephyrJung/QiaoYiCommunity/config"
	"github.com/ZephyrJung/QiaoYiCommunity/handler/api"
	"github.com/ZephyrJung/QiaoYiCommunity/middleware"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/jwt"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/oss"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/redis"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/sms"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/wechat"
	"github.com/ZephyrJung/QiaoYiCommunity/repository/dao"
	"github.com/ZephyrJung/QiaoYiCommunity/router"
	"github.com/ZephyrJung/QiaoYiCommunity/service"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("load config failed: %v", err))
	}

	gin.SetMode(cfg.Server.Mode)

	if err := utils.InitSnowflake(1); err != nil {
		panic(fmt.Sprintf("init snowflake failed: %v", err))
	}

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("connect database failed: %v", err))
	}

	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		panic(fmt.Sprintf("connect redis failed: %v", err))
	}
	defer redisClient.Close()

	userDAO := dao.NewUserDAO(db)
	categoryDAO := dao.NewCategoryDAO(db)
	itemDAO := dao.NewItemDAO(db)
	feedbackDAO := dao.NewFeedbackDAO(db)
	postDAO := dao.NewPostDAO(db)
	commentDAO := dao.NewCommentDAO(db)
	likeDAO := dao.NewLikeDAO(db)

	jwtMgr := jwt.NewManager(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.AccessExpireHours, cfg.JWT.RefreshExpireDays)
	var smsProvider sms.Provider = sms.NewMockProvider()
	wechatMP := wechat.NewMiniProgram(cfg.WeChat.AppID, cfg.WeChat.AppSecret)

	var storage oss.Storage
	if cfg.OSS.Type == "local" {
		storage = oss.NewLocalStorage("uploads", cfg.OSS.Domain)
	} else {
		storage = oss.NewLocalStorage("uploads", cfg.OSS.Domain)
	}

	authService := service.NewAuthService(userDAO, jwtMgr, redisClient, smsProvider, wechatMP)
	userService := service.NewUserService(userDAO)
	itemService := service.NewItemService(itemDAO, categoryDAO)
	feedbackService := service.NewFeedbackService(feedbackDAO)
	postService := service.NewPostService(postDAO, categoryDAO)
	commentService := service.NewCommentService(commentDAO)
	likeService := service.NewLikeService(likeDAO)
	uploadService := service.NewUploadService(storage)

	authHandler := api.NewAuthHandler(authService)
	userHandler := api.NewUserHandler(userService)
	itemHandler := api.NewItemHandler(itemService)
	feedbackHandler := api.NewFeedbackHandler(feedbackService)
	postHandler := api.NewPostHandler(postService)
	commentHandler := api.NewCommentHandler(commentService)
	likeHandler := api.NewLikeHandler(likeService)
	uploadHandler := api.NewUploadHandler(uploadService)
	categoryHandler := api.NewCategoryHandler(categoryDAO)

	rateLimiter := middleware.NewRateLimiter(redisClient, 100, 60)

	deps := &router.Dependencies{
		AuthHandler:     authHandler,
		UserHandler:     userHandler,
		ItemHandler:     itemHandler,
		FeedbackHandler: feedbackHandler,
		PostHandler:     postHandler,
		CommentHandler:  commentHandler,
		LikeHandler:     likeHandler,
		UploadHandler:   uploadHandler,
		CategoryHandler: categoryHandler,
		JWTMgr:          jwtMgr,
		RateLimiter:     rateLimiter,
	}

	r := router.SetupRouter(deps)

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	go func() {
		fmt.Printf("Server starting on %s\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("server listen error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("server shutdown error: %v\n", err)
	}
	fmt.Println("Server exited")
}
