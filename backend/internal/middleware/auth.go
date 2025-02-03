package middleware

import (
	"gym/pkg/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWT 认证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "未登录或非法访问",
			})
			c.Abort()
			return
		}

		// 检查Authorization格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "请求头Authorization格式无效",
			})
			c.Abort()
			return
		}

		// 解析token
		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "无效的token",
			})
			c.Abort()
			return
		}

		// 将用户信息保存到上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}
