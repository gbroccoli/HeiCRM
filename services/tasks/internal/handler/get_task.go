package handler

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetTask returns a single task with details
func (h *Handler) GetTask(c *gin.Context) {
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

	query := `
		SELECT t.id, t.author_id, t.assignee_id, t.room_id, t.task_type, t.description,
		       t.priority, t.status, t.created_at, t.updated_at,
		       author.name, assignee.name, r.room_number, r.building_id
		FROM tasks t
		JOIN users author ON author.id = t.author_id
		LEFT JOIN users assignee ON assignee.id = t.assignee_id
		JOIN rooms r ON r.id = t.room_id
		WHERE t.id = $1
	`

	var task models.TaskWithDetails
	err = h.DB.QueryRow(query, taskID).Scan(
		&task.ID, &task.AuthorID, &task.AssigneeID, &task.RoomID, &task.TaskType, &task.Description,
		&task.Priority, &task.Status, &task.CreatedAt, &task.UpdatedAt,
		&task.AuthorName, &task.AssigneeName, &task.RoomNumber, &task.BuildingID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Заявка не найдена")
		return
	}
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Check access: regular users can only view their own tasks or tasks assigned to them
	if userRole == 0 {
		isAuthor := task.AuthorID == userID
		isAssignee := task.AssigneeID != nil && *task.AssigneeID == userID
		if !isAuthor && !isAssignee {
			response.Forbidden(c, "Нет доступа к этой заявке")
			return
		}
	}

	response.SuccessOK(c, task)
}
