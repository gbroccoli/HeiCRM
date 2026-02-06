package handler

import (
	"database/sql"
	"errors"
)

// getUserIDByEmail returns the user ID for the given email
func getUserIDByEmail(db *sql.DB, email string) (uint64, error) {
	var userID uint64
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("user not found")
	}
	return userID, err
}

// taskExists checks if a task with the given ID exists
func taskExists(db *sql.DB, taskID uint64) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)", taskID).Scan(&exists)
	return exists, err
}

// roomExists checks if a room with the given ID exists
func roomExists(db *sql.DB, roomID uint64) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)", roomID).Scan(&exists)
	return exists, err
}

// userExists checks if a user with the given ID exists
func userExists(db *sql.DB, userID uint64) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
	return exists, err
}

// isTaskAuthor checks if the user is the author of the task
func isTaskAuthor(db *sql.DB, taskID, userID uint64) (bool, error) {
	var authorID uint64
	err := db.QueryRow("SELECT author_id FROM tasks WHERE id = $1", taskID).Scan(&authorID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return authorID == userID, nil
}

// isTaskAssignee checks if the user is the assignee of the task
func isTaskAssignee(db *sql.DB, taskID, userID uint64) (bool, error) {
	var assigneeID *uint64
	err := db.QueryRow("SELECT assignee_id FROM tasks WHERE id = $1", taskID).Scan(&assigneeID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if assigneeID == nil {
		return false, nil
	}
	return *assigneeID == userID, nil
}

// getTaskStatus returns the current status of a task
func getTaskStatus(db *sql.DB, taskID uint64) (string, error) {
	var status string
	err := db.QueryRow("SELECT status FROM tasks WHERE id = $1", taskID).Scan(&status)
	return status, err
}
