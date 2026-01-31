package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gbroccoli/HeiCRM/services/users/internal/middleware"
	"github.com/gin-gonic/gin"
)

// GetUser returns a user by ID. Accessible by admin, staff, or the user themselves.
func (h *Handler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	// Check access: admin/manager can view anyone, regular users can only view themselves
	role, roleOk := middleware.GetUserRole(c)
	if !roleOk {
		response.Unauthorized(c, "role not found in context")
		return
	}

	isAdminOrManager := role == middleware.RoleAdmin || role == middleware.RoleManager

	if !isAdminOrManager {
		// Regular user — check if requesting own profile
		email, exists := c.Get("email")
		if !exists {
			response.Unauthorized(c, "email not found in context")
			return
		}
		currentUserID, err := getUserIDByEmail(h.DB, email.(string))
		if err != nil {
			response.Unauthorized(c, "user not found")
			return
		}
		if currentUserID != userID {
			response.Forbidden(c, "insufficient permissions")
			return
		}
	}

	user, err := models.GetUserWithProfile(h.DB, userID)
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
