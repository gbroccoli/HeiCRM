package nats

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/gbroccoli/HeiCRM/pkg/events"
	"github.com/gbroccoli/HeiCRM/services/notification/internal/email"
	"github.com/nats-io/nats.go"
)

// Subscriber handles NATS subscriptions for notification events
type Subscriber struct {
	NC     *nats.Conn
	DB     *sql.DB
	Sender *email.Sender
}

// NewSubscriber creates a new notification NATS subscriber
func NewSubscriber(nc *nats.Conn, db *sql.DB, sender *email.Sender) *Subscriber {
	return &Subscriber{
		NC:     nc,
		DB:     db,
		Sender: sender,
	}
}

// Subscribe sets up all NATS subscriptions
func (s *Subscriber) Subscribe() error {
	if _, err := s.NC.Subscribe(events.SubjectUserRegistered, s.handleUserRegistered); err != nil {
		return err
	}
	log.Println("subscribed to", events.SubjectUserRegistered)

	if _, err := s.NC.Subscribe(events.SubjectTaskAssigned, s.handleTaskAssigned); err != nil {
		return err
	}
	log.Println("subscribed to", events.SubjectTaskAssigned)

	if _, err := s.NC.Subscribe(events.SubjectTaskStatusChanged, s.handleTaskStatusChanged); err != nil {
		return err
	}
	log.Println("subscribed to", events.SubjectTaskStatusChanged)

	return nil
}

// handleUserRegistered sends welcome email with credentials
func (s *Subscriber) handleUserRegistered(msg *nats.Msg) {
	var event events.UserRegisteredEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("failed to unmarshal user.registered event: %v", err)
		return
	}

	log.Printf("received user.registered for user_id=%d email=%s", event.UserID, event.Email)

	if event.Password == "" {
		log.Printf("skipping welcome email for user_id=%d: no password in event", event.UserID)
		return
	}

	body := email.RenderWelcome(event.Name, event.Email, event.Password)
	if err := s.Sender.Send(event.Email, "Добро пожаловать в HeiCRM", body); err != nil {
		log.Printf("failed to send welcome email to %s: %v", event.Email, err)
		return
	}

	log.Printf("sent welcome email to %s", event.Email)
}

// handleTaskAssigned sends notification to the assignee
func (s *Subscriber) handleTaskAssigned(msg *nats.Msg) {
	var event events.TaskAssignedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("failed to unmarshal task.assigned event: %v", err)
		return
	}

	log.Printf("received task.assigned for task_id=%d assignee_id=%d", event.TaskID, event.AssigneeID)

	// Look up assignee email
	assigneeEmail, err := getUserEmailByID(s.DB, event.AssigneeID)
	if err != nil {
		log.Printf("failed to get assignee email for user_id=%d: %v", event.AssigneeID, err)
		return
	}

	body := email.RenderTaskAssigned(event.TaskType, event.Priority, event.Description)
	if err := s.Sender.Send(assigneeEmail, "Вам назначена заявка — HeiCRM", body); err != nil {
		log.Printf("failed to send task assigned email to %s: %v", assigneeEmail, err)
		return
	}

	log.Printf("sent task assigned email to %s", assigneeEmail)
}

// handleTaskStatusChanged sends notification to the task author
func (s *Subscriber) handleTaskStatusChanged(msg *nats.Msg) {
	var event events.TaskStatusChangedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("failed to unmarshal task.status_changed event: %v", err)
		return
	}

	log.Printf("received task.status_changed for task_id=%d: %s -> %s", event.TaskID, event.PreviousStatus, event.NewStatus)

	// Notify the author about status change
	authorEmail, err := getUserEmailByID(s.DB, event.AuthorID)
	if err != nil {
		log.Printf("failed to get author email for user_id=%d: %v", event.AuthorID, err)
		return
	}

	body := email.RenderTaskStatusChanged(event.TaskType, event.PreviousStatus, event.NewStatus, event.Description)
	if err := s.Sender.Send(authorEmail, "Статус заявки изменён — HeiCRM", body); err != nil {
		log.Printf("failed to send status changed email to %s: %v", authorEmail, err)
		return
	}

	log.Printf("sent status changed email to %s", authorEmail)
}

// getUserEmailByID returns the email for a given user ID
func getUserEmailByID(db *sql.DB, userID uint64) (string, error) {
	var email string
	err := db.QueryRow("SELECT email FROM users WHERE id = $1", userID).Scan(&email)
	return email, err
}
