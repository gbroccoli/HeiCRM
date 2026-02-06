package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// ListRooms returns paginated list of rooms for a building
func (h *Handler) ListRooms(c *gin.Context) {
	buildingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID здания")
		return
	}

	exists, err := buildingExists(h.DB, buildingID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !exists {
		response.NotFoundError(c, "Здание не найдено")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	statusFilter := c.Query("status")
	floorFilter := c.Query("floor")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// Build query with filters
	countQuery := "SELECT COUNT(*) FROM rooms WHERE building_id = $1"
	args := []interface{}{buildingID}
	argIdx := 2

	if statusFilter != "" {
		countQuery += " AND status = $" + strconv.Itoa(argIdx)
		args = append(args, statusFilter)
		argIdx++
	}
	if floorFilter != "" {
		countQuery += " AND floor = $" + strconv.Itoa(argIdx)
		floor, _ := strconv.Atoi(floorFilter)
		args = append(args, floor)
		argIdx++
	}

	var total int64
	err = h.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	query := `
		SELECT r.id, r.building_id, r.room_number, r.floor, r.capacity,
		       r.room_type, r.status, r.created_at, r.updated_at,
		       COUNT(res.id) AS occupancy
		FROM rooms r
		LEFT JOIN residents res ON res.room_id = r.id AND res.move_out_date IS NULL
		WHERE r.building_id = $1
	`

	queryArgs := []interface{}{buildingID}
	queryArgIdx := 2

	if statusFilter != "" {
		query += " AND r.status = $" + strconv.Itoa(queryArgIdx)
		queryArgs = append(queryArgs, statusFilter)
		queryArgIdx++
	}
	if floorFilter != "" {
		query += " AND r.floor = $" + strconv.Itoa(queryArgIdx)
		floor, _ := strconv.Atoi(floorFilter)
		queryArgs = append(queryArgs, floor)
		queryArgIdx++
	}

	query += `
		GROUP BY r.id
		ORDER BY r.floor ASC, r.room_number ASC
		LIMIT $` + strconv.Itoa(queryArgIdx) + ` OFFSET $` + strconv.Itoa(queryArgIdx+1)
	queryArgs = append(queryArgs, pageSize, offset)

	rows, err := h.DB.Query(query, queryArgs...)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	defer rows.Close()

	type RoomWithOccupancy struct {
		models.Room
		Occupancy int `json:"occupancy"`
	}

	var rooms []RoomWithOccupancy
	for rows.Next() {
		var r RoomWithOccupancy
		err := rows.Scan(
			&r.ID, &r.BuildingID, &r.RoomNumber, &r.Floor, &r.Capacity,
			&r.RoomType, &r.Status, &r.CreatedAt, &r.UpdatedAt,
			&r.Occupancy,
		)
		if err != nil {
			response.DatabaseErrorResponse(c, err)
			return
		}
		rooms = append(rooms, r)
	}

	if err := rows.Err(); err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"code": response.OK,
		"data": models.ListResponse{
			Items: rooms,
			Pagination: models.PaginationResponse{
				Page:       page,
				Limit:      pageSize,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	})
}
