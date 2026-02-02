package midleware

import (
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	RoleUser    uint = 0
	RoleAdmin   uint = 1
	RoleManager uint = 2
)

func RoleMiddleware(allowedRoles ...uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			response.Unauthorized(c, "Роль не найдена в контексте")
			c.Abort()
			return
		}

		userRole, ok := roleVal.(uint)
		if !ok {
			response.Unauthorized(c, "Некорректный тип роли")
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Недостаточно прав")
		c.Abort()
	}
}
