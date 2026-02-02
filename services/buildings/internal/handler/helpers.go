package handler

import (
	"database/sql"
	"errors"
)

// buildingExists checks if a building with the given ID exists
func buildingExists(db *sql.DB, buildingID uint64) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM buildings WHERE id = $1)", buildingID).Scan(&exists)
	return exists, err
}

// roomExists checks if a room with the given ID belongs to the given building
func roomExists(db *sql.DB, buildingID, roomID uint64) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1 AND building_id = $2)", roomID, buildingID).Scan(&exists)
	return exists, err
}

// getRoomCurrentOccupancy returns how many residents are currently in a room
func getRoomCurrentOccupancy(db *sql.DB, roomID uint64) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user_profiles WHERE room_id = $1", roomID).Scan(&count)
	return count, err
}

// getRoomCapacity returns the capacity of a room
func getRoomCapacity(db *sql.DB, roomID uint64) (int, error) {
	var capacity int
	err := db.QueryRow("SELECT capacity FROM rooms WHERE id = $1", roomID).Scan(&capacity)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return capacity, err
}
