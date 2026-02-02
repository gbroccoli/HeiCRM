package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// getUserByEmail retrieves user from database by email
func getUserByEmail(db *sql.DB, email string) (*models.User, error) {
	user := &models.User{}

	query := `SELECT id, name, avatar, email, role_id, tg_send, created_at, updated_at FROM users WHERE email = $1`
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Avatar,
		&user.Email,
		&user.RoleID,
		&user.TgSend,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *Handler) MeUsers(c *gin.Context) {
	// Get email from JWT claims (set by AuthMiddleware)
	email, exists := c.Get("email")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": response.AuthRequired,
			"msg":  "Не авторизован",
		})
		return
	}

	emailStr, ok := email.(string)
	if !ok || emailStr == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": response.AuthRequired,
			"msg":  "Некорректный email в токене",
		})
		return
	}

	// Get user from database by email
	user, err := getUserByEmail(h.DB, emailStr)
	if err != nil {
		log.Printf("Failed to get user by email: %v", err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code": response.NotFound,
			"msg":  "Пользователь не найден",
		})
		return
	}

	// Get role name from database
	roleName, err := user.GetRole(h.DB)
	if err != nil {
		log.Printf("Failed to get user role: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": response.InternalError,
			"msg":  "Не удалось получить роль",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   response.OK,
		"id":     user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"avatar": user.Avatar,
		"role":   roleName,
	})
}
