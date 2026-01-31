package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// UserProfile represents a user's extended profile
type UserProfile struct {
	ID          uint64     `json:"id" db:"id"`
	UserID      uint64     `json:"user_id" db:"user_id"`
	FirstName   *string    `json:"first_name,omitempty" db:"first_name"`
	LastName    *string    `json:"last_name,omitempty" db:"last_name"`
	MiddleName  *string    `json:"middle_name,omitempty" db:"middle_name"`
	Phone       *string    `json:"phone,omitempty" db:"phone"`
	StudentID   *string    `json:"student_id,omitempty" db:"student_id"`
	RoomID      *uint64    `json:"room_id,omitempty" db:"room_id"`
	AvatarURL   *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" db:"date_of_birth"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// UpdateProfileRequest is the request body for updating own profile
type UpdateProfileRequest struct {
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	MiddleName *string `json:"middle_name"`
	Phone      *string `json:"phone"`
}

// AdminUpdateUserRequest is the request body for admin updating a user
type AdminUpdateUserRequest struct {
	Name        *string `json:"name"`
	RoleID      *uint   `json:"role_id"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	MiddleName  *string `json:"middle_name"`
	Phone       *string `json:"phone"`
	StudentID   *string `json:"student_id"`
	RoomID      *uint64 `json:"room_id"`
	AvatarURL   *string `json:"avatar_url"`
	DateOfBirth *string `json:"date_of_birth"`
}

// UserWithProfile combines user and profile data for API responses
type UserWithProfile struct {
	ID          uint64     `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Avatar      *string    `json:"avatar,omitempty"`
	RoleID      uint       `json:"role_id"`
	RoleName    *string    `json:"role_name,omitempty"`
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	MiddleName  *string    `json:"middle_name,omitempty"`
	Phone       *string    `json:"phone,omitempty"`
	StudentID   *string    `json:"student_id,omitempty"`
	RoomID      *uint64    `json:"room_id,omitempty"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// ActivityLog represents a user activity log entry
type ActivityLog struct {
	ID        uint64          `json:"id" db:"id"`
	UserID    uint64          `json:"user_id" db:"user_id"`
	Action    string          `json:"action" db:"action"`
	Details   json.RawMessage `json:"details,omitempty" db:"details"`
	IPAddress *string         `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent *string         `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

// ListResponse wraps paginated list results
type ListResponse struct {
	Items      interface{}        `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

// GetUserWithProfile fetches user joined with profile by user ID
func GetUserWithProfile(db *sql.DB, userID uint64) (*UserWithProfile, error) {
	query := `
		SELECT u.id, u.name, u.email, u.avatar, u.role_id, r.name as role_name,
		       p.first_name, p.last_name, p.middle_name, p.phone,
		       p.student_id, p.room_id, p.avatar_url, p.date_of_birth,
		       u.created_at, u.updated_at
		FROM users u
		LEFT JOIN user_profiles p ON u.id = p.user_id
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`

	uwp := &UserWithProfile{}
	err := db.QueryRow(query, userID).Scan(
		&uwp.ID, &uwp.Name, &uwp.Email, &uwp.Avatar, &uwp.RoleID, &uwp.RoleName,
		&uwp.FirstName, &uwp.LastName, &uwp.MiddleName, &uwp.Phone,
		&uwp.StudentID, &uwp.RoomID, &uwp.AvatarURL, &uwp.DateOfBirth,
		&uwp.CreatedAt, &uwp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return uwp, nil
}
