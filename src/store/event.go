package store

import (
	"time"
)

// Event is a created, updated or deleted event
type Event struct {
	At time.Time `json:"at" firestore:"at"`
	By string    `json:"by" firestore:"by"`
}

// NewEvent returns a new Event
func NewEvent(userID string) *Event {
	return &Event{
		At: time.Now().UTC(),
		By: userID,
	}
}
