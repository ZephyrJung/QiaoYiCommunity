package router

import (
	_ "github.com/ZephyrJung/QiaoYiCommunity/docs"
	"github.com/ZephyrJung/QiaoYiCommunity/handler/api"
	"github.com/ZephyrJung/QiaoYiCommunity/middleware"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/jwt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Dependencies struct {
	AuthHandler     *api.AuthHandler
	UserHandler     *api.UserHandler
	ItemHandler     *api.ItemHandler
	FeedbackHandler *api.FeedbackHandler
	PostHandler     *api.PostHandler
	CommentHandler  *api.CommentHandler
	LikeHandler     *api.LikeHandler
	UploadHandler   *api.UploadHandler
	CategoryHandler *api.CategoryHandler
	JWTMgr          *jwt.Manager
	RateLimiter     *middleware.RateLimiter
}

func SetupRouter(deps *Dependencies) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(deps.RateLimiter.Middleware())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	public := r.Group("/api/v1")
	{
		public.POST("/auth/wechat/login", deps.AuthHandler.WeChatLogin)
		public.POST("/auth/phone/send-code", deps.AuthHandler.PhoneSendCode)
		public.POST("/auth/phone/login", deps.AuthHandler.PhoneLogin)
		public.POST("/auth/refresh", deps.AuthHandler.RefreshToken)

		public.GET("/categories", deps.CategoryHandler.List)

		public.GET("/items", deps.ItemHandler.List)
		public.GET("/items/:id", deps.ItemHandler.GetByID)

		public.GET("/posts", deps.PostHandler.List)
		public.GET("/posts/:id", deps.PostHandler.GetByID)

		public.GET("/comments/post/:post_id", deps.CommentHandler.ListByPostID)
	}

	auth := public.Group("")
	auth.Use(middleware.JWTAuth(deps.JWTMgr))
	{
		auth.POST("/auth/logout", deps.AuthHandler.Logout)

		auth.GET("/user/profile", deps.UserHandler.GetProfile)
		auth.PUT("/user/profile", deps.UserHandler.UpdateProfile)

		auth.POST("/items", deps.ItemHandler.Create)
		auth.PUT("/items/:id", deps.ItemHandler.Update)
		auth.DELETE("/items/:id", deps.ItemHandler.Delete)
		auth.PATCH("/items/:id/sold", deps.ItemHandler.MarkSold)

		auth.POST("/feedbacks", deps.FeedbackHandler.Create)
		auth.GET("/feedbacks", deps.FeedbackHandler.List)
		auth.GET("/feedbacks/:id", deps.FeedbackHandler.GetByID)
		auth.PUT("/feedbacks/:id", deps.FeedbackHandler.Update)
		auth.DELETE("/feedbacks/:id", deps.FeedbackHandler.Delete)
		auth.POST("/feedbacks/:id/replies", deps.FeedbackHandler.CreateReply)

		auth.POST("/posts", deps.PostHandler.Create)
		auth.PUT("/posts/:id", deps.PostHandler.Update)
		auth.DELETE("/posts/:id", deps.PostHandler.Delete)

		auth.POST("/comments", deps.CommentHandler.Create)
		auth.DELETE("/comments/:id", deps.CommentHandler.Delete)

		auth.POST("/likes/toggle", deps.LikeHandler.Toggle)
		auth.GET("/likes/status", deps.LikeHandler.Status)

		auth.POST("/upload/images", deps.UploadHandler.UploadImages)
	}

	return r
}
