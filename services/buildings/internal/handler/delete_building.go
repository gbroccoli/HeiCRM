package handler

import (
	"log"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// DeleteBuilding deletes a building and its rooms (CASCADE)
func (h *Handler) DeleteBuilding(c *gin.Context) {
	idParam := c.Param("id")
	buildingID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID здания")
		return
	}

	result, err := h.DB.Exec("DELETE FROM buildings WHERE id = $1", buildingID)
	if err != nil {
		log.Printf("failed to delete building %d: %v", buildingID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	if rowsAffected == 0 {
		response.NotFoundError(c, "Здание не найдено")
		return
	}

	response.SuccessDeleted(c, nil)
}
