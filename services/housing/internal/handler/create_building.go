package handler

import (
	"log"
	"strings"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// CreateBuilding creates a new building
func (h *Handler) CreateBuilding(c *gin.Context) {
	var req models.CreateBuildingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректные данные", err)
		return
	}

	var b models.Building
	err := h.DB.QueryRow(
		`INSERT INTO buildings (address, floors, description) VALUES ($1, $2, $3)
		 RETURNING id, address, floors, description, created_at, updated_at`,
		req.Address, req.Floors, req.Description,
	).Scan(&b.ID, &b.Address, &b.Floors, &b.Description, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			response.AlreadyExistsError(c, "Здание с таким адресом уже существует")
			return
		}
		log.Printf("failed to create building: %v", err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessCreated(c, b)
}
