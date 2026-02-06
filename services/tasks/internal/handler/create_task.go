package handler

import (
	"log"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// CreateTask creates a new task
func (h *Handler) CreateTask(c *gin.Context) {
	// Get author from JWT
	email, exists := c.Get("email")
	if !exists {
		response.Unauthorized(c, "Пользователь не авторизован")
		return
	}

	authorID, err := getUserIDByEmail(h.DB, email.(string))
	if err != nil {
		response.NotFoundError(c, "Пользователь не найден")
		return
	}

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректные данные", err)
		return
	}

	// Check if room exists
	exists2, err := roomExists(h.DB, req.RoomID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !exists2 {
		response.NotFoundError(c, "Комната не найдена")
		return
	}

	var task models.Task
	err = h.DB.QueryRow(
		`INSERT INTO tasks (author_id, room_id, task_type, description, priority, status)
		 VALUES ($1, $2, $3, $4, $5, 'new')
		 RETURNING id, author_id, assignee_id, room_id, task_type, description, priority, status, created_at, updated_at`,
		authorID, req.RoomID, req.TaskType, req.Description, req.Priority,
	).Scan(&task.ID, &task.AuthorID, &task.AssigneeID, &task.RoomID,
		&task.TaskType, &task.Description, &task.Priority, &task.Status,
		&task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		log.Printf("failed to create task: %v", err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Add initial history entry
	_, err = h.DB.Exec(
		`INSERT INTO task_history (task_id, previous_status, new_status, changed_by, comment)
		 VALUES ($1, '', 'new', $2, 'Заявка создана')`,
		task.ID, authorID,
	)
	if err != nil {
		log.Printf("failed to create task history: %v", err)
	}

	response.SuccessCreated(c, task)
}
