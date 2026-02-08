package handler

import (
	"strconv"
	"time"

	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// AttachmentWithUploader extends TaskAttachment with uploader name
type AttachmentWithUploader struct {
	ID           uint64    `json:"id"`
	TaskID       uint64    `json:"task_id"`
	FileName     string    `json:"file_name"`
	FileSize     *int64    `json:"file_size,omitempty"`
	UploadedBy   uint64    `json:"uploaded_by"`
	UploadedAt   time.Time `json:"uploaded_at"`
	UploaderName string    `json:"uploader_name"`
}

// ListAttachments returns all attachments for a task
func (h *Handler) ListAttachments(c *gin.Context) {
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

	rows, err := h.DB.Query(
		`SELECT a.id, a.task_id, a.file_name, a.file_size, a.uploaded_by, a.uploaded_at, u.name
		 FROM task_attachments a
		 JOIN users u ON u.id = a.uploaded_by
		 WHERE a.task_id = $1
		 ORDER BY a.uploaded_at DESC`,
		taskID,
	)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	defer rows.Close()

	var attachments []AttachmentWithUploader
	for rows.Next() {
		var a AttachmentWithUploader
		if err := rows.Scan(&a.ID, &a.TaskID, &a.FileName, &a.FileSize,
			&a.UploadedBy, &a.UploadedAt, &a.UploaderName); err != nil {
			response.DatabaseErrorResponse(c, err)
			return
		}
		attachments = append(attachments, a)
	}

	if err := rows.Err(); err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessOK(c, attachments)
}
