package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// User represents a user entity in the database
type User struct {
	ID        uint64     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Avatar    *string    `json:"avatar,omitempty" db:"avatar"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"-" db:"password"` // Never expose password in JSON
	RoleID    uint       `json:"role_id" db:"role_id"`
	TgSend    *bool      `json:"tg_send" db:"tg_send"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type Role struct {
	ID          uint       `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

func (u *User) GetMe(db *sql.DB) (*User, error) {
	user := User{}

	query := `SELECT id, name, avatar, email, role_id FROM users WHERE id = $1`
	err := db.QueryRow(query, u.ID).Scan(&user.ID, &user.Name, &user.Avatar, &user.Email, &user.RoleID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) GetRole(db *sql.DB) (*string, error) {

	role := &Role{}

	query := `SELECT name FROM roles WHERE id = $1`
	err := db.QueryRow(query, &u.RoleID).Scan(&role.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("role not found")
	}

	if err != nil {
		return nil, err
	}

	return &role.Name, nil
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the user registration request payload
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
	RoleID   uint64 `json:"role_id" binding:"required"`
	TgSend   *bool  `json:"tg_send"`
}

// IsTgSend validates TgSend parameter and returns false if nil
func (r *RegisterRequest) IsTgSend() bool {
	if r.TgSend == nil {
		return false
	}
	return *r.TgSend
}
