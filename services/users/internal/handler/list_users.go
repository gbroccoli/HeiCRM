package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// ListUsers returns paginated list of users with profiles
func (h *Handler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// Count total
	var total int64
	err := h.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Query users with profiles
	query := `
		SELECT u.id, u.name, u.email, u.avatar, u.role_id, r.name as role_name,
		       p.first_name, p.last_name, p.middle_name, p.phone,
		       p.student_id, p.room_id, p.avatar_url, p.date_of_birth,
		       u.created_at, u.updated_at
		FROM users u
		LEFT JOIN user_profiles p ON u.id = p.user_id
		LEFT JOIN roles r ON u.role_id = r.id
		ORDER BY u.id ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := h.DB.Query(query, pageSize, offset)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	defer rows.Close()

	var users []models.UserWithProfile
	for rows.Next() {
		var u models.UserWithProfile
		err := rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.Avatar, &u.RoleID, &u.RoleName,
			&u.FirstName, &u.LastName, &u.MiddleName, &u.Phone,
			&u.StudentID, &u.RoomID, &u.AvatarURL, &u.DateOfBirth,
			&u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			response.DatabaseErrorResponse(c, err)
			return
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"code": response.OK,
		"data": models.ListResponse{
			Items: users,
			Pagination: models.PaginationResponse{
				Page:       page,
				Limit:      pageSize,
				Total:      total,
				TotalPages: totalPages,
			},
		},
	})
}
