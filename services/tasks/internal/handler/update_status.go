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

// UpdateTaskStatus changes the status of a task
func (h *Handler) UpdateTaskStatus(c *gin.Context) {
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

	var req models.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректные данные", err)
		return
	}

	// Get current task
	var currentStatus string
	var authorID uint64
	var assigneeID *uint64
	err = h.DB.QueryRow("SELECT status, author_id, assignee_id FROM tasks WHERE id = $1", taskID).
		Scan(&currentStatus, &authorID, &assigneeID)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Заявка не найдена")
		return
	}
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Validate status transition
	if !isValidStatusTransition(currentStatus, req.Status, userRole, userID, authorID, assigneeID) {
		response.BadRequest(c, "Недопустимый переход статуса")
		return
	}

	// Check permissions
	if userRole == 0 {
		// Regular users can only close their own tasks
		if req.Status == "closed" && authorID != userID {
			response.Forbidden(c, "Только автор может закрыть заявку")
			return
		}
		// And can't change other statuses
		if req.Status != "closed" {
			response.Forbidden(c, "Недостаточно прав для изменения статуса")
			return
		}
	}

	// Update status
	var task models.Task
	err = h.DB.QueryRow(
		`UPDATE tasks SET status = $1, updated_at = NOW() WHERE id = $2
		 RETURNING id, author_id, assignee_id, room_id, task_type, description, priority, status, created_at, updated_at`,
		req.Status, taskID,
	).Scan(&task.ID, &task.AuthorID, &task.AssigneeID, &task.RoomID,
		&task.TaskType, &task.Description, &task.Priority, &task.Status,
		&task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		log.Printf("failed to update task status: %v", err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Add history entry
	_, err = h.DB.Exec(
		`INSERT INTO task_history (task_id, previous_status, new_status, changed_by, comment)
		 VALUES ($1, $2, $3, $4, $5)`,
		taskID, currentStatus, req.Status, userID, req.Comment,
	)
	if err != nil {
		log.Printf("failed to create task history: %v", err)
	}

	response.SuccessUpdated(c, task)
}

// isValidStatusTransition checks if the status change is allowed
func isValidStatusTransition(from, to string, role uint, userID, authorID uint64, assigneeID *uint64) bool {
	// Admin can do anything
	if role == 1 {
		return true
	}

	// Define allowed transitions
	transitions := map[string][]string{
		"new":         {"assigned", "in_progress", "closed"},
		"assigned":    {"in_progress", "new", "closed"},
		"in_progress": {"completed", "assigned", "closed"},
		"completed":   {"closed", "in_progress"},
		"closed":      {}, // closed is final
	}

	allowed, ok := transitions[from]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}

	return false
}
