package midleware

import (
	"net/http"

	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	_ "github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/gbroccoli/HeiCRM/services/auth/internal/tools"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(j *jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		token, err := tools.ExtractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		claims, err := j.VerifyAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

func RefreshTokenMiddleware(j *jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}
		if refreshToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		token := refreshToken
		claims, err := j.VerifyRefreshToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Set("email", claims.Email)
		c.Set("userID", claims.ID)
		c.Next()
	}
}
