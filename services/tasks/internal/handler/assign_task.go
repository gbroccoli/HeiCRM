package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// AssignTask assigns a task to a user
func (h *Handler) AssignTask(c *gin.Context) {
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

	var req models.AssignTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректные данные", err)
		return
	}

	// Check task exists
	var currentStatus string
	err = h.DB.QueryRow("SELECT status FROM tasks WHERE id = $1", taskID).Scan(&currentStatus)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Заявка не найдена")
		return
	}
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Can't assign completed or closed tasks
	if currentStatus == "completed" || currentStatus == "closed" {
		response.BadRequest(c, "Нельзя назначить исполнителя на завершённую заявку")
		return
	}

	// Check assignee exists
	assigneeExists, err := userExists(h.DB, req.AssigneeID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !assigneeExists {
		response.NotFoundError(c, "Исполнитель не найден")
		return
	}

	// Update task
	newStatus := currentStatus
	if currentStatus == "new" {
		newStatus = "assigned"
	}

	var task models.Task
	err = h.DB.QueryRow(
		`UPDATE tasks SET assignee_id = $1, status = $2, updated_at = NOW() WHERE id = $3
		 RETURNING id, author_id, assignee_id, room_id, task_type, description, priority, status, created_at, updated_at`,
		req.AssigneeID, newStatus, taskID,
	).Scan(&task.ID, &task.AuthorID, &task.AssigneeID, &task.RoomID,
		&task.TaskType, &task.Description, &task.Priority, &task.Status,
		&task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		log.Printf("failed to assign task: %v", err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Add history entry if status changed
	if currentStatus != newStatus {
		_, err = h.DB.Exec(
			`INSERT INTO task_history (task_id, previous_status, new_status, changed_by, comment)
			 VALUES ($1, $2, $3, $4, $5)`,
			taskID, currentStatus, newStatus, userID, "Назначен исполнитель",
		)
		if err != nil {
			log.Printf("failed to create task history: %v", err)
		}
	}

	// Publish task.assigned event
	if h.NC != nil {
		event := struct {
			TaskID      uint64  `json:"task_id"`
			AssigneeID  uint64  `json:"assignee_id"`
			AuthorID    uint64  `json:"author_id"`
			TaskType    string  `json:"task_type"`
			Description string  `json:"description"`
			Priority    string  `json:"priority"`
		}{taskID, req.AssigneeID, task.AuthorID, task.TaskType, task.Description, task.Priority}
		data, err := json.Marshal(event)
		if err != nil {
			log.Printf("failed to marshal task.assigned event: %v", err)
		} else if err := h.NC.Publish("task.assigned", data); err != nil {
			log.Printf("failed to publish task.assigned event: %v", err)
		}
	}

	response.SuccessUpdated(c, task)
}

// TakeTask allows a user to assign the task to themselves
func (h *Handler) TakeTask(c *gin.Context) {
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

	// Check task exists and is not already assigned
	var currentStatus string
	var currentAssignee *uint64
	err = h.DB.QueryRow("SELECT status, assignee_id FROM tasks WHERE id = $1", taskID).Scan(&currentStatus, &currentAssignee)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Заявка не найдена")
		return
	}
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	if currentStatus == "completed" || currentStatus == "closed" {
		response.BadRequest(c, "Нельзя взять завершённую заявку")
		return
	}

	if currentAssignee != nil {
		response.BadRequest(c, "Заявка уже назначена другому исполнителю")
		return
	}

	// Assign to self
	newStatus := "assigned"
	if currentStatus == "assigned" {
		newStatus = currentStatus
	}

	var task models.Task
	err = h.DB.QueryRow(
		`UPDATE tasks SET assignee_id = $1, status = $2, updated_at = NOW() WHERE id = $3
		 RETURNING id, author_id, assignee_id, room_id, task_type, description, priority, status, created_at, updated_at`,
		userID, newStatus, taskID,
	).Scan(&task.ID, &task.AuthorID, &task.AssigneeID, &task.RoomID,
		&task.TaskType, &task.Description, &task.Priority, &task.Status,
		&task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		log.Printf("failed to take task: %v", err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Add history entry
	if currentStatus != newStatus {
		_, err = h.DB.Exec(
			`INSERT INTO task_history (task_id, previous_status, new_status, changed_by, comment)
			 VALUES ($1, $2, $3, $4, $5)`,
			taskID, currentStatus, newStatus, userID, "Взял заявку в работу",
		)
		if err != nil {
			log.Printf("failed to create task history: %v", err)
		}
	}

	// Publish task.assigned event
	if h.NC != nil {
		event := struct {
			TaskID      uint64 `json:"task_id"`
			AssigneeID  uint64 `json:"assignee_id"`
			AuthorID    uint64 `json:"author_id"`
			TaskType    string `json:"task_type"`
			Description string `json:"description"`
			Priority    string `json:"priority"`
		}{taskID, userID, task.AuthorID, task.TaskType, task.Description, task.Priority}
		data, err := json.Marshal(event)
		if err != nil {
			log.Printf("failed to marshal task.assigned event: %v", err)
		} else if err := h.NC.Publish("task.assigned", data); err != nil {
			log.Printf("failed to publish task.assigned event: %v", err)
		}
	}

	response.SuccessUpdated(c, task)
}
