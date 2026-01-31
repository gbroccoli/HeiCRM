package middleware

import (
	"strings"

	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies JWT access token and sets user info in context
func AuthMiddleware(j *jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			response.InvalidTokenError(c)
			c.Abort()
			return
		}

		claims, err := j.VerifyAccessToken(token)
		if err != nil {
			response.InvalidTokenError(c)
			c.Abort()
			return
		}

		c.Set("email", claims.Subject)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RoleMiddleware checks if the user has one of the allowed roles
func RoleMiddleware(allowedRoles ...uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			response.Unauthorized(c, "role not found in context")
			c.Abort()
			return
		}

		userRole, ok := roleVal.(uint)
		if !ok {
			response.Unauthorized(c, "invalid role type")
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "insufficient permissions")
		c.Abort()
	}
}

// Role constants matching the roles table
const (
	RoleUser    uint = 0
	RoleAdmin   uint = 1
	RoleManager uint = 2
)

// GetUserRole extracts user role from gin context
func GetUserRole(c *gin.Context) (uint, bool) {
	role, exists := c.Get("role")
	if !exists {
		return 0, false
	}
	r, ok := role.(uint)
	return r, ok
}
