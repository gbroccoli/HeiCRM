package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   uint   `json:"role_id"`
	TgSend   *bool  `json:"tg_send"`
}

func GetUser(db *sql.DB, email, password string) (*User, error) {
	user := &User{}

	query := `SELECT id, name, email, password, role_id, tg_send FROM users WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&user.Id, &user.Name, user.Email, user.Password, &user.RoleID, &user.TgSend)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not fount")
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *Handler) Login(c *gin.Context) {

	var userParams LoginRequest
	if err := c.ShouldBindJSON(&userParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса."})
		return
	}

	// logic no database users
	user, err := GetUser(h.DB, userParams.Email, userParams.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkPassword := h.PasswordManager.CheckHash(user.Password, userParams.Password)
	if !checkPassword {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": 0401,
			"msg":  "Invalid login or password",
		})
		return
	}

	tokenAccess, err := h.JWT.GenerateAccessToken(userParams.Email, 1, "access")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Fatal(err.Error())
		return
	}

	refreshToken, err := h.JWT.GenerateRefreshToken(userParams.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Fatal(err.Error())
		return
	}

	expires := 30 * 24 * time.Hour
	c.SetCookie(
		"refresh",
		refreshToken,
		int(expires),
		"/auth/refresh",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "success login",
		"token":   tokenAccess,
	})
}
