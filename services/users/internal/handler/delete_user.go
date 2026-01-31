package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// DeleteUser allows admin to delete a user
func (h *Handler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	// Prevent self-deletion
	email, exists := c.Get("email")
	if exists {
		currentUserID, err := getUserIDByEmail(h.DB, email.(string))
		if err == nil && currentUserID == userID {
			response.Forbidden(c, "cannot delete yourself")
			return
		}
	}

	result, err := h.DB.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		log.Printf("failed to delete user %d: %v", userID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	if rowsAffected == 0 {
		response.NotFoundError(c, "user not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    response.Deleted,
		"message": "user deleted",
	})
}
