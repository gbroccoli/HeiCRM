package handler

import (
	"log"
	"strconv"
	"strings"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// CreateRoom creates a new room in a building
func (h *Handler) CreateRoom(c *gin.Context) {
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

	var req models.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректные данные", err)
		return
	}

	// Validate floor does not exceed building floors
	var floors int
	h.DB.QueryRow("SELECT floors FROM buildings WHERE id = $1", buildingID).Scan(&floors)
	if req.Floor > floors {
		response.BadRequest(c, "Этаж не может превышать количество этажей в здании")
		return
	}

	var room models.Room
	err = h.DB.QueryRow(
		`INSERT INTO rooms (building_id, room_number, floor, capacity, room_type, status)
		 VALUES ($1, $2, $3, $4, $5, 'free')
		 RETURNING id, building_id, room_number, floor, capacity, room_type, status, created_at, updated_at`,
		buildingID, req.RoomNumber, req.Floor, req.Capacity, req.RoomType,
	).Scan(&room.ID, &room.BuildingID, &room.RoomNumber, &room.Floor,
		&room.Capacity, &room.RoomType, &room.Status, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			response.AlreadyExistsError(c, "Комната с таким номером уже существует в этом здании")
			return
		}
		log.Printf("failed to create room: %v", err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessCreated(c, room)
}
