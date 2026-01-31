package handler

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetActivity returns activity log for a specific user (admin only)
func (h *Handler) GetActivity(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize

	// Count total
	var total int64
	err = h.DB.QueryRow("SELECT COUNT(*) FROM user_activity_log WHERE user_id = $1", userID).Scan(&total)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Query activity log
	query := `
		SELECT id, user_id, action, details, ip_address, user_agent, created_at
		FROM user_activity_log
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := h.DB.Query(query, userID, pageSize, offset)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	defer rows.Close()

	var logs []models.ActivityLog
	for rows.Next() {
		var l models.ActivityLog
		var details []byte
		err := rows.Scan(&l.ID, &l.UserID, &l.Action, &details, &l.IPAddress, &l.UserAgent, &l.CreatedAt)
		if err != nil {
			response.DatabaseErrorResponse(c, err)
			return
		}
		if details != nil {
			l.Details = json.RawMessage(details)
		}
		logs = append(logs, l)
	}

	if err := rows.Err(); err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"code": response.OK,
		"data": models.ListResponse{
			Items: logs,
			Pagination: models.PaginationResponse{
				Page:       page,
				Limit:      pageSize,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	})
}
