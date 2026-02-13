package events

// NATS subjects
const (
	SubjectUserRegistered    = "user.registered"
	SubjectProfileUpdated    = "user.profile_updated"
	SubjectUserDeactivated   = "user.deactivated"
	SubjectTaskAssigned      = "task.assigned"
	SubjectTaskStatusChanged = "task.status_changed"
)

// UserRegisteredEvent is published by Auth Service when a new user is registered.
type UserRegisteredEvent struct {
	UserID   uint64 `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// ProfileUpdatedEvent is published by User Service when a user profile is updated.
type ProfileUpdatedEvent struct {
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
}

// UserDeactivatedEvent is published by User Service when a user is deactivated/deleted.
type UserDeactivatedEvent struct {
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
}

// TaskAssignedEvent is published by Tasks Service when a task is assigned.
type TaskAssignedEvent struct {
	TaskID      uint64 `json:"task_id"`
	AssigneeID  uint64 `json:"assignee_id"`
	AuthorID    uint64 `json:"author_id"`
	TaskType    string `json:"task_type"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// TaskStatusChangedEvent is published by Tasks Service when a task status changes.
type TaskStatusChangedEvent struct {
	TaskID         uint64  `json:"task_id"`
	AuthorID       uint64  `json:"author_id"`
	AssigneeID     *uint64 `json:"assignee_id,omitempty"`
	PreviousStatus string  `json:"previous_status"`
	NewStatus      string  `json:"new_status"`
	ChangedBy      uint64  `json:"changed_by"`
	TaskType       string  `json:"task_type"`
	Description    string  `json:"description"`
}
