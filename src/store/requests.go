package store

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// RequestsCollection is the name of the collection
const RequestsCollection = "requests"

// Request represents a model in the database
type Request struct {
	ID        string            `json:"id" firestore:"id"`
	Name      string            `json:"name" firestore:"name"`
	FolderID  string            `json:"folder_id" firestore:"folder_id"`
	ProjectID string            `json:"project_id" firestore:"project_id"`
	Type      string            `json:"type" firestore:"type"`
	URL       string            `json:"url" firestore:"url"`
	Headers   map[string]string `json:"headers" firestore:"headers"`
	Body      string            `json:"body" firestore:"body"`
	Created   *Event            `json:"created" firestore:"created"`
	Updated   *Event            `json:"updated" firestore:"updated"`
	Deleted   *Event            `json:"deleted" firestore:"deleted"`
}

// NewRequestID generates a UUID for requests
func (s *Store) NewRequestID() string {
	return "req-" + uuid.New().String()
}

// CreateRequest creates a new request
func (s *Store) CreateRequest(ctx context.Context, userID string, request *Request) error {
	request.Created = NewEvent(userID)
	_, err := s.Client.Collection(RequestsCollection).Doc(request.ID).Set(ctx, request)
	return err
}

// DeleteRequest deletes an existing request
func (s *Store) DeleteRequest(ctx context.Context, userID string, request *Request) error {
	request.Deleted = NewEvent(userID)
	_, err := s.Client.Collection(RequestsCollection).Doc(request.ID).Set(ctx, request)
	return err
}

// UpdateRequest updates an existing request
func (s *Store) UpdateRequest(ctx context.Context, userID string, request *Request) error {
	request.Updated = NewEvent(userID)
	_, err := s.Client.Collection(RequestsCollection).Doc(request.ID).Set(ctx, request)
	return err
}

// GetRequestByID gets a Request by id
func (s *Store) GetRequestByID(ctx context.Context, id string) (*Request, error) {
	iter := s.Client.Collection(RequestsCollection).Where("id", "==", id).Limit(1).Documents(ctx)

	snapshot, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	request := &Request{}
	snapshot.DataTo(request)

	if request.Deleted != nil {
		return nil, nil
	}

	return request, nil
}
