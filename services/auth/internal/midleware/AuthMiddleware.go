package midleware

import (
	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	_ "github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(j *jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		}

		jwtJson, err := j.Verify(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		}

		if "access" != jwtJson.Type || jwtJson.Type != "refresh" {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		}

		c.Next()
	}
}
