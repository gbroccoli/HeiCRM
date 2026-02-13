package nats

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/gbroccoli/HeiCRM/pkg/events"
	"github.com/nats-io/nats.go"
)

// Subscriber handles NATS message subscriptions
type Subscriber struct {
	NC *nats.Conn
	DB *sql.DB
}

// NewSubscriber creates a new NATS subscriber
func NewSubscriber(nc *nats.Conn, db *sql.DB) *Subscriber {
	return &Subscriber{
		NC: nc,
		DB: db,
	}
}

// Subscribe sets up all NATS subscriptions
func (s *Subscriber) Subscribe() error {
	_, err := s.NC.Subscribe(events.SubjectUserRegistered, s.handleUserRegistered)
	if err != nil {
		return err
	}

	log.Println("subscribed to user.registered")
	return nil
}

// handleUserRegistered creates a profile entry when a new user is registered
func (s *Subscriber) handleUserRegistered(msg *nats.Msg) {
	var event events.UserRegisteredEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("failed to unmarshal user.registered event: %v", err)
		return
	}

	log.Printf("received user.registered event for user_id=%d email=%s", event.UserID, event.Email)

	query := `
		INSERT INTO user_profiles (user_id, created_at, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO NOTHING
	`

	now := time.Now()
	_, err := s.DB.Exec(query, event.UserID, now, now)
	if err != nil {
		log.Printf("failed to create profile for user %d: %v", event.UserID, err)
		return
	}

	log.Printf("created profile for user %d", event.UserID)
}

// PublishProfileUpdated publishes a profile updated event
func PublishProfileUpdated(nc *nats.Conn, userID uint64, email string) {
	event := events.ProfileUpdatedEvent{
		UserID: userID,
		Email:  email,
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal profile updated event: %v", err)
		return
	}

	if err := nc.Publish(events.SubjectProfileUpdated, data); err != nil {
		log.Printf("failed to publish user.profile_updated: %v", err)
	}
}

// PublishUserDeactivated publishes a user deactivated event
func PublishUserDeactivated(nc *nats.Conn, userID uint64, email string) {
	event := events.UserDeactivatedEvent{
		UserID: userID,
		Email:  email,
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal user deactivated event: %v", err)
		return
	}

	if err := nc.Publish(events.SubjectUserDeactivated, data); err != nil {
		log.Printf("failed to publish user.deactivated: %v", err)
	}
}
