package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/ZephyrJung/QiaoYiCommunity/pkg/response"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				fmt.Printf("[PANIC] %v\n%s\n", err, stack)
				response.InternalError(c, "internal server error")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
