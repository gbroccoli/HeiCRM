package handler

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// DeleteAttachment removes a task attachment (only uploader or admin)
func (h *Handler) DeleteAttachment(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID заявки")
		return
	}

	attachmentID, err := strconv.ParseUint(c.Param("attachmentId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID вложения")
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

	// Get attachment
	var uploadedBy uint64
	var filePath string
	err = h.DB.QueryRow(
		"SELECT uploaded_by, file_path FROM task_attachments WHERE id = $1 AND task_id = $2",
		attachmentID, taskID,
	).Scan(&uploadedBy, &filePath)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Вложение не найдено")
		return
	}
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Only uploader or admin can delete
	if uploadedBy != userID && userRole != 1 {
		response.Forbidden(c, "Только загрузивший или администратор может удалить вложение")
		return
	}

	// Delete file from disk
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		log.Printf("failed to remove file %s: %v", filePath, err)
	}

	// Delete from DB
	_, err = h.DB.Exec("DELETE FROM task_attachments WHERE id = $1", attachmentID)
	if err != nil {
		log.Printf("failed to delete attachment %d: %v", attachmentID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessDeleted(c, nil)
}
