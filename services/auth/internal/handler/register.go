package handler

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/gbroccoli/HeiCRM/pkg/events"
	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

func createUser(db *sql.DB, user *models.RegisterRequest) (uint64, error) {
	query := `INSERT INTO users (name, email, password, role_id, tg_send) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	var id uint64
	err := db.QueryRow(query, user.Name, user.Email, user.Password, user.RoleID, user.TgSend).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (h *Handler) Register(c *gin.Context) {

	var candidate models.RegisterRequest
	if err := c.ShouldBindJSON(&candidate); err != nil {
		response.ValidationError(c, "Некорректные данные запроса")
		return
	}

	// create user
	userDraft := &models.RegisterRequest{
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
		response.BadRequestError(c, "Не удалось создать хеш пароля", err)
		log.Printf("Failed to generate password hash: %v", err)
		return
	}

	userDraft.Password = string(hash)

	userID, err := createUser(h.DB, userDraft)
	if err != nil {
		response.InternalErrorResponse(c, "Не удалось создать пользователя", err)
		log.Printf("Failed to create user: %v", err)
		return
	}

	log.Printf("User created: id=%d email=%s password=%s", userID, candidate.Email, password)

	// publish user.registered event to NATS
	event := events.UserRegisteredEvent{
		UserID:   userID,
		Email:    candidate.Email,
		Name:     candidate.Name,
		Password: password,
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal user.registered event: %v", err)
	} else if err := h.NC.Publish(events.SubjectUserRegistered, data); err != nil {
		log.Printf("Failed to publish user.registered event: %v", err)
	} else {
		log.Printf("Published user.registered event for user_id=%d", userID)
	}

	response.SuccessCreated(c, gin.H{
		"user_id": userID,
		"email":   candidate.Email,
		"name":    candidate.Name,
	})
}
