package handler

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetBuilding returns a single building with room/resident stats
func (h *Handler) GetBuilding(c *gin.Context) {
	idParam := c.Param("id")
	buildingID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID здания")
		return
	}

	query := `
		SELECT b.id, b.address, b.floors, b.description, b.created_at, b.updated_at,
		       COALESCE(rs.room_count, 0),
		       COALESCE(rs.resident_count, 0)
		FROM buildings b
		LEFT JOIN (
			SELECT r.building_id,
			       COUNT(DISTINCT r.id) AS room_count,
			       COUNT(res.id) AS resident_count
			FROM rooms r
			LEFT JOIN residents res ON res.room_id = r.id AND res.move_out_date IS NULL
			GROUP BY r.building_id
		) rs ON rs.building_id = b.id
		WHERE b.id = $1
	`

	var b models.BuildingWithStats
	err = h.DB.QueryRow(query, buildingID).Scan(
		&b.ID, &b.Address, &b.Floors, &b.Description, &b.CreatedAt, &b.UpdatedAt,
		&b.RoomCount, &b.ResidentCount,
	)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Здание не найдено")
		return
	}
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessOK(c, b)
}
