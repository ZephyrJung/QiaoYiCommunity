package middleware

import (
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/jwt"
	"github.com/ZephyrJung/QiaoYiCommunity/pkg/response"
)

const (
	ContextUserIDKey  = "user_id"
	ContextRoleKey    = "role"
	ContextTokenIDKey = "token_id"
)

func JWTAuth(jwtMgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "invalid authorization format")
			c.Abort()
			return
		}

		claims, err := jwtMgr.ParseAccessToken(parts[1])
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextRoleKey, claims.Role)
		c.Set(ContextTokenIDKey, claims.TokenID)
		c.Next()
	}
}

func RequireRole(minRole int8) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextRoleKey)
		if !exists {
			response.Forbidden(c, "role not found")
			c.Abort()
			return
		}
		if r, ok := role.(int8); !ok || r < minRole {
			response.Forbidden(c, "insufficient permissions")
			c.Abort()
			return
		}
		c.Next()
	}
}
