package models

import "time"

// BaseModel contains common fields for database entities
type BaseModel struct {
	ID        uint64     `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// PaginationRequest represents common pagination parameters
type PaginationRequest struct {
	Page  int `json:"page" form:"page" binding:"omitempty,min=1"`
	Limit int `json:"limit" form:"limit" binding:"omitempty,min=1,max=100"`
}

// DefaultPagination returns default pagination values
func (p *PaginationRequest) DefaultPagination() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 20
	}
}

// Offset calculates the offset for SQL queries
func (p *PaginationRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}

// PaginationResponse represents pagination metadata in API responses
type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}
