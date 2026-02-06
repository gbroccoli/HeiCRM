package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// ListTasks returns paginated list of tasks with filters
func (h *Handler) ListTasks(c *gin.Context) {
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	statusFilter := c.Query("status")
	priorityFilter := c.Query("priority")
	assigneeFilter := c.Query("assignee")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// Build base query
	baseQuery := `
		FROM tasks t
		JOIN users author ON author.id = t.author_id
		LEFT JOIN users assignee ON assignee.id = t.assignee_id
		JOIN rooms r ON r.id = t.room_id
		WHERE 1=1
	`

	args := []interface{}{}
	argIdx := 1

	// For regular users, only show their own tasks
	if userRole == 0 {
		baseQuery += " AND t.author_id = $" + strconv.Itoa(argIdx)
		args = append(args, userID)
		argIdx++
	}

	// Apply filters
	if statusFilter != "" {
		baseQuery += " AND t.status = $" + strconv.Itoa(argIdx)
		args = append(args, statusFilter)
		argIdx++
	}
	if priorityFilter != "" {
		baseQuery += " AND t.priority = $" + strconv.Itoa(argIdx)
		args = append(args, priorityFilter)
		argIdx++
	}
	if assigneeFilter == "me" {
		baseQuery += " AND t.assignee_id = $" + strconv.Itoa(argIdx)
		args = append(args, userID)
		argIdx++
	} else if assigneeFilter == "unassigned" {
		baseQuery += " AND t.assignee_id IS NULL"
	} else if assigneeFilter != "" {
		assigneeID, err := strconv.ParseUint(assigneeFilter, 10, 64)
		if err == nil {
			baseQuery += " AND t.assignee_id = $" + strconv.Itoa(argIdx)
			args = append(args, assigneeID)
			argIdx++
		}
	}

	// Count total
	var total int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	err = h.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Fetch tasks
	selectQuery := `
		SELECT t.id, t.author_id, t.assignee_id, t.room_id, t.task_type, t.description,
		       t.priority, t.status, t.created_at, t.updated_at,
		       author.name, assignee.name, r.room_number, r.building_id
	` + baseQuery + `
		ORDER BY
			CASE t.priority
				WHEN 'critical' THEN 1
				WHEN 'high' THEN 2
				WHEN 'medium' THEN 3
				WHEN 'low' THEN 4
			END,
			t.created_at DESC
		LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)

	args = append(args, pageSize, offset)

	rows, err := h.DB.Query(selectQuery, args...)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	defer rows.Close()

	var tasks []models.TaskWithDetails
	for rows.Next() {
		var t models.TaskWithDetails
		err := rows.Scan(
			&t.ID, &t.AuthorID, &t.AssigneeID, &t.RoomID, &t.TaskType, &t.Description,
			&t.Priority, &t.Status, &t.CreatedAt, &t.UpdatedAt,
			&t.AuthorName, &t.AssigneeName, &t.RoomNumber, &t.BuildingID,
		)
		if err != nil {
			response.DatabaseErrorResponse(c, err)
			return
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"code": response.OK,
		"data": models.ListResponse{
			Items: tasks,
			Pagination: models.PaginationResponse{
				Page:       page,
				Limit:      pageSize,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	})
}
