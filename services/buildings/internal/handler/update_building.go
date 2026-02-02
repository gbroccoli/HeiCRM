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

// UpdateBuilding updates an existing building
func (h *Handler) UpdateBuilding(c *gin.Context) {
	idParam := c.Param("id")
	buildingID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID здания")
		return
	}

	var req models.UpdateBuildingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректные данные", err)
		return
	}

	// Check building exists
	exists, err := buildingExists(h.DB, buildingID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !exists {
		response.NotFoundError(c, "Здание не найдено")
		return
	}

	// Build dynamic update
	query := "UPDATE buildings SET updated_at = NOW()"
	args := []interface{}{}
	argIdx := 1

	if req.Name != nil {
		query += ", name = $" + strconv.Itoa(argIdx)
		args = append(args, *req.Name)
		argIdx++
	}
	if req.Address != nil {
		query += ", address = $" + strconv.Itoa(argIdx)
		args = append(args, *req.Address)
		argIdx++
	}
	if req.Floors != nil {
		query += ", floors = $" + strconv.Itoa(argIdx)
		args = append(args, *req.Floors)
		argIdx++
	}

	query += " WHERE id = $" + strconv.Itoa(argIdx) +
		" RETURNING id, name, address, floors, created_at, updated_at"
	args = append(args, buildingID)

	var b models.Building
	err = h.DB.QueryRow(query, args...).Scan(
		&b.ID, &b.Name, &b.Address, &b.Floors, &b.CreatedAt, &b.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Здание не найдено")
		return
	}
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			response.AlreadyExistsError(c, "Здание с таким названием уже существует")
			return
		}
		log.Printf("failed to update building %d: %v", buildingID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	response.SuccessUpdated(c, b)
}
