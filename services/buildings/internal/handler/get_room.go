package handler

import (
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetRoom returns a room with its residents
func (h *Handler) GetRoom(c *gin.Context) {
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

	// Fetch room
	var room models.RoomWithResidents
	err = h.DB.QueryRow(
		`SELECT id, building_id, number, floor, capacity, created_at, updated_at
		 FROM rooms WHERE id = $1 AND building_id = $2`,
		roomID, buildingID,
	).Scan(&room.ID, &room.BuildingID, &room.Number, &room.Floor,
		&room.Capacity, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		response.NotFoundError(c, "Комната не найдена")
		return
	}

	// Fetch residents
	rows, err := h.DB.Query(
		`SELECT u.id, u.name, u.email, up.first_name, up.last_name
		 FROM user_profiles up
		 JOIN users u ON u.id = up.user_id
		 WHERE up.room_id = $1
		 ORDER BY u.name ASC`,
		roomID,
	)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	defer rows.Close()

	room.Residents = []models.RoomResident{}
	for rows.Next() {
		var r models.RoomResident
		if err := rows.Scan(&r.UserID, &r.Name, &r.Email, &r.FirstName, &r.LastName); err != nil {
			response.DatabaseErrorResponse(c, err)
			return
		}
		room.Residents = append(room.Residents, r)
	}

	if err := rows.Err(); err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessOK(c, room)
}
