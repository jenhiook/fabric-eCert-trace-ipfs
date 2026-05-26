package middleware

import (
	"backend/pkg"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "请求未携带token，无权限访问1",
			})
			c.Abort()
			return
		}
		mc, err := pkg.ParseToken(authHeader)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "无效的token",
				"data": err.Error(),
			})
			c.Abort()
			return
		}
		// 将当前请求的用户信息保存到上下文中
		c.Set("userID", mc.UserID)
		c.Set("userType", mc.UserType)   // 新增：保存用户角色
		c.Next()
	}
}

// RoleCheck 角色检查中间件，必须放在 JWTAuthMiddleware 之后使用
func RoleCheck(allowedRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("userType")
		if !exists {
			c.JSON(403, gin.H{
				"code": 403,
				"msg":  "无法获取用户角色",
			})
			c.Abort()
			return
		}
		if userType != allowedRole {
			c.JSON(403, gin.H{
				"code": 403,
				"msg":  "权限不足，需要 " + allowedRole + " 角色",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
