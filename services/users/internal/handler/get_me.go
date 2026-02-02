package handler

import (
	"database/sql"
	"net/http"

	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetMe returns the current authenticated user's profile
func (h *Handler) GetMe(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		response.Unauthorized(c, "Email не найден в контексте")
		return
	}

	emailStr, ok := email.(string)
	if !ok || emailStr == "" {
		response.Unauthorized(c, "Некорректный email в токене")
		return
	}

	userID, err := getUserIDByEmail(h.DB, emailStr)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFoundError(c, "Пользователь не найден")
			return
		}
		response.DatabaseErrorResponse(c, err)
		return
	}

	user, err := getUserWithProfile(h.DB, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFoundError(c, "Пользователь не найден")
			return
		}
		response.DatabaseErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": response.OK,
		"data": user,
	})
}
