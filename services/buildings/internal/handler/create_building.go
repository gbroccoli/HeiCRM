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
		`INSERT INTO buildings (name, address, floors) VALUES ($1, $2, $3)
		 RETURNING id, name, address, floors, created_at, updated_at`,
		req.Name, req.Address, req.Floors,
	).Scan(&b.ID, &b.Name, &b.Address, &b.Floors, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			response.AlreadyExistsError(c, "Здание с таким названием уже существует")
			return
		}
		log.Printf("failed to create building: %v", err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessCreated(c, b)
}
