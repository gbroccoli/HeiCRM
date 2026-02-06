package handler

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// UpdateRoom updates an existing room
func (h *Handler) UpdateRoom(c *gin.Context) {
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

	var req models.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректные данные", err)
		return
	}

	// Validate floor if provided
	if req.Floor != nil {
		var floors int
		h.DB.QueryRow("SELECT floors FROM buildings WHERE id = $1", buildingID).Scan(&floors)
		if *req.Floor > floors {
			response.BadRequest(c, "Этаж не может превышать количество этажей в здании")
			return
		}
	}

	// Build dynamic update
	query := "UPDATE rooms SET updated_at = NOW()"
	args := []interface{}{}
	argIdx := 1

	if req.RoomNumber != nil {
		query += ", room_number = $" + strconv.Itoa(argIdx)
		args = append(args, *req.RoomNumber)
		argIdx++
	}
	if req.Floor != nil {
		query += ", floor = $" + strconv.Itoa(argIdx)
		args = append(args, *req.Floor)
		argIdx++
	}
	if req.Capacity != nil {
		query += ", capacity = $" + strconv.Itoa(argIdx)
		args = append(args, *req.Capacity)
		argIdx++
	}
	if req.RoomType != nil {
		query += ", room_type = $" + strconv.Itoa(argIdx)
		args = append(args, *req.RoomType)
		argIdx++
	}
	if req.Status != nil {
		query += ", status = $" + strconv.Itoa(argIdx)
		args = append(args, *req.Status)
		argIdx++
	}

	query += " WHERE id = $" + strconv.Itoa(argIdx) + " AND building_id = $" + strconv.Itoa(argIdx+1) +
		" RETURNING id, building_id, room_number, floor, capacity, room_type, status, created_at, updated_at"
	args = append(args, roomID, buildingID)

	var room models.Room
	err = h.DB.QueryRow(query, args...).Scan(
		&room.ID, &room.BuildingID, &room.RoomNumber, &room.Floor,
		&room.Capacity, &room.RoomType, &room.Status, &room.CreatedAt, &room.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Комната не найдена")
		return
	}
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			response.AlreadyExistsError(c, "Комната с таким номером уже существует в этом здании")
			return
		}
		log.Printf("failed to update room %d: %v", roomID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessUpdated(c, room)
}
