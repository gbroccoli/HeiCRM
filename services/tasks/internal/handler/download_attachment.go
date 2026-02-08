package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// DownloadAttachment serves a task attachment file for download
func (h *Handler) DownloadAttachment(c *gin.Context) {
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

	// Get attachment from DB
	var fileName, filePath string
	err = h.DB.QueryRow(
		"SELECT file_name, file_path FROM task_attachments WHERE id = $1 AND task_id = $2",
		attachmentID, taskID,
	).Scan(&fileName, &filePath)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Вложение не найдено")
		return
	}
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Check file exists on disk
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.NotFoundError(c, "Файл не найден на диске")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.File(filePath)
}
