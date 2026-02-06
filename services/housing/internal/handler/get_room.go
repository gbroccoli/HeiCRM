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
		`SELECT id, building_id, room_number, floor, capacity, room_type, status, created_at, updated_at
		 FROM rooms WHERE id = $1 AND building_id = $2`,
		roomID, buildingID,
	).Scan(&room.ID, &room.BuildingID, &room.RoomNumber, &room.Floor,
		&room.Capacity, &room.RoomType, &room.Status, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		response.NotFoundError(c, "Комната не найдена")
		return
	}

	// Fetch residents (only active - without move_out_date)
	rows, err := h.DB.Query(
		`SELECT id, room_id, full_name, birth_date, passport_series, passport_number,
		        email, phone, move_in_date, move_out_date, created_at, updated_at
		 FROM residents
		 WHERE room_id = $1 AND move_out_date IS NULL
		 ORDER BY full_name ASC`,
		roomID,
	)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	defer rows.Close()

	room.Residents = []models.Resident{}
	for rows.Next() {
		var r models.Resident
		if err := rows.Scan(&r.ID, &r.RoomID, &r.FullName, &r.BirthDate,
			&r.PassportSeries, &r.PassportNumber, &r.Email, &r.Phone,
			&r.MoveInDate, &r.MoveOutDate, &r.CreatedAt, &r.UpdatedAt); err != nil {
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
