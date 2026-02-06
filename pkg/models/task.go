package models

import "time"

// TaskPriority constants
const (
	TaskPriorityLow      = "low"
	TaskPriorityMedium   = "medium"
	TaskPriorityHigh     = "high"
	TaskPriorityCritical = "critical"
)

// TaskStatus constants
const (
	TaskStatusNew        = "new"
	TaskStatusAssigned   = "assigned"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
	TaskStatusClosed     = "closed"
)

// Task represents a task/request entity
type Task struct {
	ID          uint64     `json:"id" db:"id"`
	AuthorID    uint64     `json:"author_id" db:"author_id"`
	AssigneeID  *uint64    `json:"assignee_id,omitempty" db:"assignee_id"`
	RoomID      uint64     `json:"room_id" db:"room_id"`
	TaskType    string     `json:"task_type" db:"task_type"`
	Description string     `json:"description" db:"description"`
	Priority    string     `json:"priority" db:"priority"`
	Status      string     `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// TaskWithDetails extends Task with author and assignee names
type TaskWithDetails struct {
	Task
	AuthorName   string  `json:"author_name"`
	AssigneeName *string `json:"assignee_name,omitempty"`
	RoomNumber   string  `json:"room_number"`
	BuildingID   uint64  `json:"building_id"`
}

// TaskHistory represents a status change history entry
type TaskHistory struct {
	ID             uint64    `json:"id" db:"id"`
	TaskID         uint64    `json:"task_id" db:"task_id"`
	PreviousStatus string    `json:"previous_status" db:"previous_status"`
	NewStatus      string    `json:"new_status" db:"new_status"`
	ChangedBy      uint64    `json:"changed_by" db:"changed_by"`
	ChangedAt      time.Time `json:"changed_at" db:"changed_at"`
	Comment        *string   `json:"comment,omitempty" db:"comment"`
}

// TaskComment represents a comment on a task
type TaskComment struct {
	ID          uint64    `json:"id" db:"id"`
	TaskID      uint64    `json:"task_id" db:"task_id"`
	AuthorID    uint64    `json:"author_id" db:"author_id"`
	CommentText string    `json:"comment_text" db:"comment_text"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// TaskCommentWithAuthor extends TaskComment with author name
type TaskCommentWithAuthor struct {
	TaskComment
	AuthorName string `json:"author_name"`
}

// TaskAttachment represents an attachment to a task
type TaskAttachment struct {
	ID         uint64    `json:"id" db:"id"`
	TaskID     uint64    `json:"task_id" db:"task_id"`
	FileName   string    `json:"file_name" db:"file_name"`
	FilePath   string    `json:"file_path" db:"file_path"`
	FileSize   *int64    `json:"file_size,omitempty" db:"file_size"`
	UploadedBy uint64    `json:"uploaded_by" db:"uploaded_by"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}

// CreateTaskRequest is the request body for creating a task
type CreateTaskRequest struct {
	RoomID      uint64 `json:"room_id" binding:"required"`
	TaskType    string `json:"task_type" binding:"required"`
	Description string `json:"description" binding:"required"`
	Priority    string `json:"priority" binding:"required,oneof=low medium high critical"`
}

// UpdateTaskRequest is the request body for updating a task
type UpdateTaskRequest struct {
	TaskType    *string `json:"task_type"`
	Description *string `json:"description"`
	Priority    *string `json:"priority" binding:"omitempty,oneof=low medium high critical"`
}

// UpdateTaskStatusRequest is the request body for changing task status
type UpdateTaskStatusRequest struct {
	Status  string  `json:"status" binding:"required,oneof=new assigned in_progress completed closed"`
	Comment *string `json:"comment"`
}

// AssignTaskRequest is the request body for assigning a task
type AssignTaskRequest struct {
	AssigneeID uint64 `json:"assignee_id" binding:"required"`
}

// CreateTaskCommentRequest is the request body for adding a comment
type CreateTaskCommentRequest struct {
	CommentText string `json:"comment_text" binding:"required"`
}
