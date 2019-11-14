package store

import (
	"apiboy/backend/src/config"

	"cloud.google.com/go/firestore"
)

// Store wraps the Firestore client
type Store struct {
	Config *config.Config
	Client *firestore.Client
}

// New returns a new Store
func New(conf *config.Config, client *firestore.Client) *Store {
	return &Store{
		Config: conf,
		Client: client,
	}
}
