package handler

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

const maxFileSize = 10 << 20 // 10MB

// UploadAttachment uploads a file attachment to a task
func (h *Handler) UploadAttachment(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID заявки")
		return
	}

	email, exists := c.Get("email")
	if !exists {
		response.Unauthorized(c, "Пользователь не авторизован")
		return
	}

	userID, err := getUserIDByEmail(h.DB, email.(string))
	if err != nil {
		response.NotFoundError(c, "Пользователь не найден")
		return
	}

	roleVal, _ := c.Get("role")
	userRole := roleVal.(uint)

	// Check task exists
	exists2, err := taskExists(h.DB, taskID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !exists2 {
		response.NotFoundError(c, "Заявка не найдена")
		return
	}

	// Check access
	canAccess, err := canAccessTask(h.DB, taskID, userID, userRole)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !canAccess {
		response.Forbidden(c, "Нет доступа к этой заявке")
		return
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "Файл не найден в запросе")
		return
	}

	if file.Size > maxFileSize {
		response.BadRequest(c, "Размер файла превышает 10MB")
		return
	}

	// Generate UUID filename
	uuid := make([]byte, 16)
	if _, err := rand.Read(uuid); err != nil {
		response.InternalErrorResponse(c, "Ошибка генерации имени файла", err)
		return
	}
	ext := filepath.Ext(file.Filename)
	uuidName := fmt.Sprintf("%x%s", uuid, ext)

	// Create directory
	dir := fmt.Sprintf("./uploads/tasks/%d", taskID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("failed to create upload dir: %v", err)
		response.InternalErrorResponse(c, "Ошибка создания директории", err)
		return
	}

	// Save file
	filePath := filepath.Join(dir, uuidName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("failed to save file: %v", err)
		response.InternalErrorResponse(c, "Ошибка сохранения файла", err)
		return
	}

	// Insert into DB
	var attachment models.TaskAttachment
	err = h.DB.QueryRow(
		`INSERT INTO task_attachments (task_id, file_name, file_path, file_size, uploaded_by)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, task_id, file_name, file_path, file_size, uploaded_by, uploaded_at`,
		taskID, file.Filename, filePath, file.Size, userID,
	).Scan(&attachment.ID, &attachment.TaskID, &attachment.FileName, &attachment.FilePath,
		&attachment.FileSize, &attachment.UploadedBy, &attachment.UploadedAt)
	if err != nil {
		log.Printf("failed to insert attachment: %v", err)
		os.Remove(filePath)
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessCreated(c, attachment)
}
