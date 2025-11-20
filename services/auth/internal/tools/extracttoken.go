package tools

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractToken(c *gin.Context) (string, error) {
	token := c.GetHeader("Authorization")

	if token == "" {
		return "", errors.New("header Authorization is empty")
	}

	parts := strings.Split(token, " ")

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header format'")
	}

	return parts[1], nil
}
