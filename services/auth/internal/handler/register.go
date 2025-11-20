package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	RoleID   uint64 `json:"role_id" binding:"required"`
	TgSend   *bool  `json:"tg_send"`
}

// IsTgSend validate params on null
func (register *RegisterRequest) IsTgSend() bool {
	if register.TgSend == nil {
		return false
	}

	return *register.TgSend
}

//tokenH, err := tools.ExtractToken(c)
//if err != nil {
//c.JSON(200, gin.H{"error": err.Error()})
//return
//}
//
//isTokenAccess, err := h.JWT.VerifyAccessToken(tokenH)
//if err != nil {
//c.JSON(200, gin.H{"error": err.Error()})
//}
//
//c.JSON(200, gin.H{"token": isTokenAccess})
//return

func CreateUser(db *sql.DB, user *RegisterRequest) (*uint64, error) {
	query := `INSERT INTO users (name, email, password, role_id, tg_send) VALUES (?, ?, ?, ?, ?)`
	res, err := db.Exec(query, user.Name, user.Email, user.Password, user.RoleID, user.TgSend)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	uid := uint64(id)
	return &uid, nil
}

func (h *Handler) Register(c *gin.Context) {

	var candidate RegisterRequest
	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса."})
		return
	}

	// create user

	// Auto-generate a 24-character password
	password := h.PasswordManager.GeneratePassword()

	// Generate the bcrypt hash to store in the database
	hash, err := h.PasswordManager.GenerateHash(password)
	if err != nil {
		log.Fatalln(hash, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"name":     candidate.Name,
			"tg_send":  candidate.IsTgSend(),
			"password": password,
			"hash":     hash,
		},
	})
}
