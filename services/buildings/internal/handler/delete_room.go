package handler

import (
	"log"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// DeleteRoom deletes a room from a building
func (h *Handler) DeleteRoom(c *gin.Context) {
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

	result, err := h.DB.Exec("DELETE FROM rooms WHERE id = $1 AND building_id = $2", roomID, buildingID)
	if err != nil {
		log.Printf("failed to delete room %d: %v", roomID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	if rowsAffected == 0 {
		response.NotFoundError(c, "Комната не найдена")
		return
	}

	response.SuccessDeleted(c, nil)
}
