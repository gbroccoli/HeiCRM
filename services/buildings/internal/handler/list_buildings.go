package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// ListBuildings returns paginated list of buildings with room/resident stats
func (h *Handler) ListBuildings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// Count total
	var total int64
	err := h.DB.QueryRow("SELECT COUNT(*) FROM buildings").Scan(&total)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	query := `
		SELECT b.id, b.name, b.address, b.floors, b.created_at, b.updated_at,
		       COALESCE(rs.room_count, 0),
		       COALESCE(rs.resident_count, 0)
		FROM buildings b
		LEFT JOIN (
			SELECT r.building_id,
			       COUNT(DISTINCT r.id) AS room_count,
			       COUNT(up.user_id) AS resident_count
			FROM rooms r
			LEFT JOIN user_profiles up ON up.room_id = r.id
			GROUP BY r.building_id
		) rs ON rs.building_id = b.id
		ORDER BY b.id ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := h.DB.Query(query, pageSize, offset)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	defer rows.Close()

	var buildings []models.BuildingWithStats
	for rows.Next() {
		var b models.BuildingWithStats
		err := rows.Scan(
			&b.ID, &b.Name, &b.Address, &b.Floors, &b.CreatedAt, &b.UpdatedAt,
			&b.RoomCount, &b.ResidentCount,
		)
		if err != nil {
			response.DatabaseErrorResponse(c, err)
			return
		}
		buildings = append(buildings, b)
	}

	if err := rows.Err(); err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"code": response.OK,
		"data": models.ListResponse{
			Items: buildings,
			Pagination: models.PaginationResponse{
				Page:       page,
				Limit:      pageSize,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	})
}
