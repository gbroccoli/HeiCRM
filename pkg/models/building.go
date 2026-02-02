package models

import "time"

// Building represents a building (dormitory) entity
type Building struct {
	ID        uint64     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Address   *string    `json:"address,omitempty" db:"address"`
	Floors    int        `json:"floors" db:"floors"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// BuildingWithStats extends Building with room/resident counts
type BuildingWithStats struct {
	Building
	RoomCount     int `json:"room_count"`
	ResidentCount int `json:"resident_count"`
}

// Room represents a room in a building
type Room struct {
	ID         uint64     `json:"id" db:"id"`
	BuildingID uint64     `json:"building_id" db:"building_id"`
	Number     string     `json:"number" db:"number"`
	Floor      int        `json:"floor" db:"floor"`
	Capacity   int        `json:"capacity" db:"capacity"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// RoomWithResidents extends Room with a list of residents
type RoomWithResidents struct {
	Room
	Residents []RoomResident `json:"residents"`
}

// RoomResident is a simplified user representation for room listings
type RoomResident struct {
	UserID    uint64  `json:"user_id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

// CreateBuildingRequest is the request body for creating a building
type CreateBuildingRequest struct {
	Name    string  `json:"name" binding:"required"`
	Address *string `json:"address"`
	Floors  int     `json:"floors" binding:"required,min=1"`
}

// UpdateBuildingRequest is the request body for updating a building
type UpdateBuildingRequest struct {
	Name    *string `json:"name"`
	Address *string `json:"address"`
	Floors  *int    `json:"floors" binding:"omitempty,min=1"`
}

// CreateRoomRequest is the request body for creating a room
type CreateRoomRequest struct {
	Number   string `json:"number" binding:"required"`
	Floor    int    `json:"floor" binding:"required,min=1"`
	Capacity int    `json:"capacity" binding:"required,min=1"`
}

// UpdateRoomRequest is the request body for updating a room
type UpdateRoomRequest struct {
	Number   *string `json:"number"`
	Floor    *int    `json:"floor" binding:"omitempty,min=1"`
	Capacity *int    `json:"capacity" binding:"omitempty,min=1"`
}

// AssignResidentRequest is the request body for assigning a resident to a room
type AssignResidentRequest struct {
	UserID uint64 `json:"user_id" binding:"required"`
}
