package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// ListResidents returns paginated list of residents for a room
func (h *Handler) ListResidents(c *gin.Context) {
	buildingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID здания")
		return
	}

	roomID, err := strconv.ParseUint(c.Param("roomId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID комнаты")
		return
	}

	exists, err := roomExists(h.DB, buildingID, roomID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !exists {
		response.NotFoundError(c, "Комната не найдена")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	includeMovedOut := c.Query("include_moved_out") == "true"

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// Count query
	countQuery := "SELECT COUNT(*) FROM residents WHERE room_id = $1"
	if !includeMovedOut {
		countQuery += " AND move_out_date IS NULL"
	}

	var total int64
	if err := h.DB.QueryRow(countQuery, roomID).Scan(&total); err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Data query
	query := `SELECT id, room_id, full_name, birth_date, passport_series, passport_number,
	                 email, phone, move_in_date, move_out_date, created_at, updated_at
	          FROM residents WHERE room_id = $1`
	if !includeMovedOut {
		query += " AND move_out_date IS NULL"
	}
	query += " ORDER BY move_in_date DESC LIMIT $2 OFFSET $3"

	rows, err := h.DB.Query(query, roomID, pageSize, offset)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	defer rows.Close()

	var residents []models.Resident
	for rows.Next() {
		var r models.Resident
		if err := rows.Scan(&r.ID, &r.RoomID, &r.FullName, &r.BirthDate,
			&r.PassportSeries, &r.PassportNumber, &r.Email, &r.Phone,
			&r.MoveInDate, &r.MoveOutDate, &r.CreatedAt, &r.UpdatedAt); err != nil {
			response.DatabaseErrorResponse(c, err)
			return
		}
		residents = append(residents, r)
	}

	if err := rows.Err(); err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"code": response.OK,
		"data": models.ListResponse{
			Items: residents,
			Pagination: models.PaginationResponse{
				Page:       page,
				Limit:      pageSize,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	})
}
