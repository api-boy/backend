package store

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// EnvironmentsCollection is the name of the collection
const EnvironmentsCollection = "environments"

// Environment represents a model in the database
type Environment struct {
	ID        string            `json:"id" firestore:"id"`
	Name      string            `json:"name" firestore:"name"`
	Variables map[string]string `json:"variables" firestore:"variables"`
	ProjectID string            `json:"project_id" firestore:"project_id"`
	Created   *Event            `json:"created" firestore:"created"`
	Updated   *Event            `json:"updated" firestore:"updated"`
	Deleted   *Event            `json:"deleted" firestore:"deleted"`
}

// NewEnvironmentID generates a UUID for environments
func (s *Store) NewEnvironmentID() string {
	return "env-" + uuid.New().String()
}

// CreateEnvironment creates a new Environment
func (s *Store) CreateEnvironment(ctx context.Context, userID string, environment *Environment) error {
	environment.Created = NewEvent(userID)
	_, err := s.Client.Collection(EnvironmentsCollection).Doc(environment.ID).Set(ctx, environment)
	return err
}

// DeleteEnvironment deletes an existing Environment
func (s *Store) DeleteEnvironment(ctx context.Context, userID string, environment *Environment) error {
	environment.Deleted = NewEvent(userID)
	_, err := s.Client.Collection(EnvironmentsCollection).Doc(environment.ID).Set(ctx, environment)
	return err
}

// UpdateEnvironment updates an existing environment
func (s *Store) UpdateEnvironment(ctx context.Context, userID string, environment *Environment) error {
	environment.Updated = NewEvent(userID)
	_, err := s.Client.Collection(EnvironmentsCollection).Doc(environment.ID).Set(ctx, environment)
	return err
}

// GetEnvironmentByID gets a Environment by id
func (s *Store) GetEnvironmentByID(ctx context.Context, id string) (*Environment, error) {
	iter := s.Client.Collection(EnvironmentsCollection).Where("id", "==", id).Limit(1).Documents(ctx)

	snapshot, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	environment := &Environment{}
	snapshot.DataTo(environment)

	if environment.Deleted != nil {
		return nil, nil
	}

	return environment, nil
}
