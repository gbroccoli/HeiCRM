package handler

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetResident returns a single resident by ID
func (h *Handler) GetResident(c *gin.Context) {
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

	residentID, err := strconv.ParseUint(c.Param("residentId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID резидента")
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

	var r models.Resident
	err = h.DB.QueryRow(
		`SELECT id, room_id, full_name, birth_date,
		        email, phone, move_in_date, move_out_date, created_at, updated_at
		 FROM residents WHERE id = $1 AND room_id = $2`,
		residentID, roomID,
	).Scan(&r.ID, &r.RoomID, &r.FullName, &r.BirthDate,
		&r.Email, &r.Phone,
		&r.MoveInDate, &r.MoveOutDate, &r.CreatedAt, &r.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Резидент не найден")
		return
	}
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessOK(c, r)
}
