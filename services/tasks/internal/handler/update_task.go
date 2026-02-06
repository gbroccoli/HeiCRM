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

// UpdateTask updates task details (type, description, priority)
func (h *Handler) UpdateTask(c *gin.Context) {
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

	// Check task exists and get author
	var authorID uint64
	var currentStatus string
	err = h.DB.QueryRow("SELECT author_id, status FROM tasks WHERE id = $1", taskID).Scan(&authorID, &currentStatus)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Заявка не найдена")
		return
	}
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Only author or admin can update task details
	if userRole == 0 && authorID != userID {
		response.Forbidden(c, "Только автор может редактировать заявку")
		return
	}

	// Can't update if task is completed or closed
	if currentStatus == "completed" || currentStatus == "closed" {
		response.BadRequest(c, "Нельзя редактировать завершённую заявку")
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректные данные", err)
		return
	}

	// Build dynamic update
	query := "UPDATE tasks SET updated_at = NOW()"
	args := []interface{}{}
	argIdx := 1

	if req.TaskType != nil {
		query += ", task_type = $" + strconv.Itoa(argIdx)
		args = append(args, *req.TaskType)
		argIdx++
	}
	if req.Description != nil {
		query += ", description = $" + strconv.Itoa(argIdx)
		args = append(args, *req.Description)
		argIdx++
	}
	if req.Priority != nil {
		query += ", priority = $" + strconv.Itoa(argIdx)
		args = append(args, *req.Priority)
		argIdx++
	}

	query += " WHERE id = $" + strconv.Itoa(argIdx) +
		" RETURNING id, author_id, assignee_id, room_id, task_type, description, priority, status, created_at, updated_at"
	args = append(args, taskID)

	var task models.Task
	err = h.DB.QueryRow(query, args...).Scan(
		&task.ID, &task.AuthorID, &task.AssigneeID, &task.RoomID,
		&task.TaskType, &task.Description, &task.Priority, &task.Status,
		&task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		log.Printf("failed to update task %d: %v", taskID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessUpdated(c, task)
}
