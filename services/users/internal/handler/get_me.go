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
		response.Unauthorized(c, "email not found in context")
		return
	}

	emailStr, ok := email.(string)
	if !ok || emailStr == "" {
		response.Unauthorized(c, "invalid email in token")
		return
	}

	userID, err := getUserIDByEmail(h.DB, emailStr)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFoundError(c, "user not found")
			return
		}
		response.DatabaseErrorResponse(c, err)
		return
	}

	user, err := getUserWithProfile(h.DB, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFoundError(c, "user not found")
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
