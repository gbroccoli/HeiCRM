package handler

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// UpdateResident updates resident data with partial updates
func (h *Handler) UpdateResident(c *gin.Context) {
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

	var req models.UpdateResidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректные данные", err)
		return
	}

	// Build dynamic update
	query := "UPDATE residents SET updated_at = NOW()"
	args := []interface{}{}
	argIdx := 1

	if req.FullName != nil {
		query += ", full_name = $" + strconv.Itoa(argIdx)
		args = append(args, *req.FullName)
		argIdx++
	}
	if req.BirthDate != nil {
		query += ", birth_date = $" + strconv.Itoa(argIdx)
		args = append(args, *req.BirthDate)
		argIdx++
	}
	if req.PassportSeries != nil {
		query += ", passport_series = $" + strconv.Itoa(argIdx)
		args = append(args, *req.PassportSeries)
		argIdx++
	}
	if req.PassportNumber != nil {
		query += ", passport_number = $" + strconv.Itoa(argIdx)
		args = append(args, *req.PassportNumber)
		argIdx++
	}
	if req.Email != nil {
		query += ", email = $" + strconv.Itoa(argIdx)
		args = append(args, *req.Email)
		argIdx++
	}
	if req.Phone != nil {
		query += ", phone = $" + strconv.Itoa(argIdx)
		args = append(args, *req.Phone)
		argIdx++
	}
	if req.MoveOutDate != nil {
		query += ", move_out_date = $" + strconv.Itoa(argIdx)
		args = append(args, *req.MoveOutDate)
		argIdx++
	}

	query += " WHERE id = $" + strconv.Itoa(argIdx) + " AND room_id = $" + strconv.Itoa(argIdx+1) +
		" RETURNING id, room_id, full_name, birth_date, passport_series, passport_number, email, phone, move_in_date, move_out_date, created_at, updated_at"
	args = append(args, residentID, roomID)

	var r models.Resident
	err = h.DB.QueryRow(query, args...).Scan(
		&r.ID, &r.RoomID, &r.FullName, &r.BirthDate,
		&r.PassportSeries, &r.PassportNumber, &r.Email, &r.Phone,
		&r.MoveInDate, &r.MoveOutDate, &r.CreatedAt, &r.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Резидент не найден")
		return
	}
	if err != nil {
		log.Printf("failed to update resident %d: %v", residentID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Update room status if move_out_date was set
	if req.MoveOutDate != nil {
		if err := updateRoomStatus(h.DB, roomID); err != nil {
			log.Printf("failed to update room status: %v", err)
		}
	}

	response.SuccessUpdated(c, r)
}
