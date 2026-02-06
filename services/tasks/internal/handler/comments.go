package handler

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// AddComment adds a comment to a task
func (h *Handler) AddComment(c *gin.Context) {
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

	var req models.CreateTaskCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректные данные", err)
		return
	}

	var comment models.TaskComment
	err = h.DB.QueryRow(
		`INSERT INTO task_comments (task_id, author_id, comment_text)
		 VALUES ($1, $2, $3)
		 RETURNING id, task_id, author_id, comment_text, created_at`,
		taskID, userID, req.CommentText,
	).Scan(&comment.ID, &comment.TaskID, &comment.AuthorID, &comment.CommentText, &comment.CreatedAt)

	if err != nil {
		log.Printf("failed to create comment: %v", err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessCreated(c, comment)
}

// GetComments returns all comments for a task
func (h *Handler) GetComments(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID заявки")
		return
	}

	// Check task exists
	exists, err := taskExists(h.DB, taskID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !exists {
		response.NotFoundError(c, "Заявка не найдена")
		return
	}

	rows, err := h.DB.Query(
		`SELECT c.id, c.task_id, c.author_id, c.comment_text, c.created_at, u.name
		 FROM task_comments c
		 JOIN users u ON u.id = c.author_id
		 WHERE c.task_id = $1
		 ORDER BY c.created_at ASC`,
		taskID,
	)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	defer rows.Close()

	var comments []models.TaskCommentWithAuthor
	for rows.Next() {
		var comment models.TaskCommentWithAuthor
		if err := rows.Scan(&comment.ID, &comment.TaskID, &comment.AuthorID, &comment.CommentText, &comment.CreatedAt, &comment.AuthorName); err != nil {
			response.DatabaseErrorResponse(c, err)
			return
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessOK(c, comments)
}

// GetHistory returns the status change history for a task
func (h *Handler) GetHistory(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID заявки")
		return
	}

	// Check task exists
	exists, err := taskExists(h.DB, taskID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !exists {
		response.NotFoundError(c, "Заявка не найдена")
		return
	}

	rows, err := h.DB.Query(
		`SELECT h.id, h.task_id, h.previous_status, h.new_status, h.changed_by, h.changed_at, h.comment, u.name
		 FROM task_history h
		 JOIN users u ON u.id = h.changed_by
		 WHERE h.task_id = $1
		 ORDER BY h.changed_at ASC`,
		taskID,
	)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	defer rows.Close()

	type HistoryWithUser struct {
		models.TaskHistory
		ChangedByName string `json:"changed_by_name"`
	}

	var history []HistoryWithUser
	for rows.Next() {
		var entry HistoryWithUser
		if err := rows.Scan(&entry.ID, &entry.TaskID, &entry.PreviousStatus, &entry.NewStatus,
			&entry.ChangedBy, &entry.ChangedAt, &entry.Comment, &entry.ChangedByName); err != nil {
			response.DatabaseErrorResponse(c, err)
			return
		}
		history = append(history, entry)
	}

	if err := rows.Err(); err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessOK(c, history)
}

// DeleteTask deletes a task (admin only)
func (h *Handler) DeleteTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID заявки")
		return
	}

	// Check task exists
	var id uint64
	err = h.DB.QueryRow("SELECT id FROM tasks WHERE id = $1", taskID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Заявка не найдена")
		return
	}
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Delete task (cascade will delete comments, history, attachments)
	_, err = h.DB.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		log.Printf("failed to delete task %d: %v", taskID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessDeleted(c, nil)
}
