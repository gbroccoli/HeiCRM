package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
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

func createUser(db *sql.DB, user *RegisterRequest) (*uint64, error) {
	query := `INSERT INTO users (name, email, password, role_id, tg_send) VALUES ($1, $2, $3, $4, $5);`
	_, err := db.Exec(query, user.Name, user.Email, user.Password, user.RoleID, user.TgSend)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) Register(c *gin.Context) {

	var candidate RegisterRequest
	if err := c.ShouldBindJSON(&candidate); err != nil {
		response.ValidationError(c, "Invalid request data")
		return
	}

	// create user
	userDraft := &RegisterRequest{
		Name:   candidate.Name,
		Email:  candidate.Email,
		RoleID: candidate.RoleID,
		TgSend: candidate.TgSend,
	}

	// Auto-generate a 24-character password
	password := h.PasswordManager.GeneratePassword()

	// Generate the bcrypt hash to store in the database
	hash, err := h.PasswordManager.GenerateHash(password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": response.InvalidData,
			"msg":  err.Error(),
		})
		log.Printf("Failed to generate password hash: %v", err)
		return
	}

	userDraft.Password = string(hash)

	_, err = createUser(h.DB, userDraft)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": response.InvalidData,
			"msg":  err.Error(),
		})
		log.Printf("Failed to create user: %v", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":     response.Created,
		"msg":      "user created",
		"password": password,
	})
}
