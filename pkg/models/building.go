package models

import "time"

// Building represents a building (dormitory) entity
type Building struct {
	ID          uint64     `json:"id" db:"id"`
	Address     string     `json:"address" db:"address"`
	Floors      int        `json:"floors" db:"floors"`
	Description *string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// BuildingWithStats extends Building with room/resident counts
type BuildingWithStats struct {
	Building
	RoomCount     int `json:"room_count"`
	ResidentCount int `json:"resident_count"`
}

// RoomType constants
const (
	RoomTypeSingle = "single"
	RoomTypeDouble = "double"
	RoomTypeBlock  = "block"
)

// RoomStatus constants
const (
	RoomStatusFree     = "free"
	RoomStatusOccupied = "occupied"
)

// Room represents a room in a building
type Room struct {
	ID         uint64     `json:"id" db:"id"`
	BuildingID uint64     `json:"building_id" db:"building_id"`
	RoomNumber string     `json:"room_number" db:"room_number"`
	Floor      int        `json:"floor" db:"floor"`
	Capacity   int        `json:"capacity" db:"capacity"`
	RoomType   string     `json:"room_type" db:"room_type"`
	Status     string     `json:"status" db:"status"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// RoomWithResidents extends Room with a list of residents
type RoomWithResidents struct {
	Room
	Residents []Resident `json:"residents"`
}

// Resident represents a person living in a room
type Resident struct {
	ID             uint64     `json:"id" db:"id"`
	RoomID         uint64     `json:"room_id" db:"room_id"`
	FullName       string     `json:"full_name" db:"full_name"`
	BirthDate      string     `json:"birth_date" db:"birth_date"`
	PassportSeries *string    `json:"passport_series,omitempty" db:"passport_series"`
	PassportNumber *string    `json:"passport_number,omitempty" db:"passport_number"`
	Email          *string    `json:"email,omitempty" db:"email"`
	Phone          *string    `json:"phone,omitempty" db:"phone"`
	MoveInDate     string     `json:"move_in_date" db:"move_in_date"`
	MoveOutDate    *string    `json:"move_out_date,omitempty" db:"move_out_date"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// CreateBuildingRequest is the request body for creating a building
type CreateBuildingRequest struct {
	Address     string  `json:"address" binding:"required"`
	Floors      int     `json:"floors" binding:"required,min=1"`
	Description *string `json:"description"`
}

// UpdateBuildingRequest is the request body for updating a building
type UpdateBuildingRequest struct {
	Address     *string `json:"address"`
	Floors      *int    `json:"floors" binding:"omitempty,min=1"`
	Description *string `json:"description"`
}

// CreateRoomRequest is the request body for creating a room
type CreateRoomRequest struct {
	RoomNumber string `json:"room_number" binding:"required"`
	Floor      int    `json:"floor" binding:"required,min=1"`
	Capacity   int    `json:"capacity" binding:"required,min=1"`
	RoomType   string `json:"room_type" binding:"required,oneof=single double block"`
}

// UpdateRoomRequest is the request body for updating a room
type UpdateRoomRequest struct {
	RoomNumber *string `json:"room_number"`
	Floor      *int    `json:"floor" binding:"omitempty,min=1"`
	Capacity   *int    `json:"capacity" binding:"omitempty,min=1"`
	RoomType   *string `json:"room_type" binding:"omitempty,oneof=single double block"`
	Status     *string `json:"status" binding:"omitempty,oneof=free occupied"`
}

// CreateResidentRequest is the request body for creating a resident
type CreateResidentRequest struct {
	FullName       string  `json:"full_name" binding:"required"`
	BirthDate      string  `json:"birth_date" binding:"required"`
	PassportSeries *string `json:"passport_series"`
	PassportNumber *string `json:"passport_number"`
	Email          *string `json:"email" binding:"omitempty,email"`
	Phone          *string `json:"phone"`
	MoveInDate     string  `json:"move_in_date" binding:"required"`
}

// UpdateResidentRequest is the request body for updating a resident
type UpdateResidentRequest struct {
	FullName       *string `json:"full_name"`
	BirthDate      *string `json:"birth_date"`
	PassportSeries *string `json:"passport_series"`
	PassportNumber *string `json:"passport_number"`
	Email          *string `json:"email" binding:"omitempty,email"`
	Phone          *string `json:"phone"`
	MoveOutDate    *string `json:"move_out_date"`
}

// TransferResidentRequest is the request body for transferring a resident to another room
type TransferResidentRequest struct {
	NewRoomID     uint64 `json:"new_room_id" binding:"required"`
	NewBuildingID uint64 `json:"new_building_id" binding:"required"`
}
