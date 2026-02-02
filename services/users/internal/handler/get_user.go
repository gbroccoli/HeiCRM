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
		response.BadRequest(c, "Некорректный ID пользователя")
		return
	}

	// Check access: admin/manager can view anyone, regular users can only view themselves
	role, roleOk := middleware.GetUserRole(c)
	if !roleOk {
		response.Unauthorized(c, "Роль не найдена в контексте")
		return
	}

	isAdminOrManager := role == middleware.RoleAdmin || role == middleware.RoleManager

	if !isAdminOrManager {
		// Regular user — check if requesting own profile
		email, exists := c.Get("email")
		if !exists {
			response.Unauthorized(c, "Email не найден в контексте")
			return
		}
		currentUserID, err := getUserIDByEmail(h.DB, email.(string))
		if err != nil {
			response.Unauthorized(c, "Пользователь не найден")
			return
		}
		if currentUserID != userID {
			response.Forbidden(c, "Недостаточно прав")
			return
		}
	}

	user, err := models.GetUserWithProfile(h.DB, userID)
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
