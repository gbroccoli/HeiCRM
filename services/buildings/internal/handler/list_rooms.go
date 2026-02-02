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

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var total int64
	err = h.DB.QueryRow("SELECT COUNT(*) FROM rooms WHERE building_id = $1", buildingID).Scan(&total)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	query := `
		SELECT r.id, r.building_id, r.number, r.floor, r.capacity,
		       r.created_at, r.updated_at,
		       COUNT(up.user_id) AS occupancy
		FROM rooms r
		LEFT JOIN user_profiles up ON up.room_id = r.id
		WHERE r.building_id = $1
		GROUP BY r.id
		ORDER BY r.floor ASC, r.number ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := h.DB.Query(query, buildingID, pageSize, offset)
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
			&r.ID, &r.BuildingID, &r.Number, &r.Floor, &r.Capacity,
			&r.CreatedAt, &r.UpdatedAt,
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
