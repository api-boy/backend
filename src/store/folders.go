package store

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// FoldersCollection is the name of the collection
const FoldersCollection = "folders"

// Folder represents a model in the database
type Folder struct {
	ID        string `json:"id" firestore:"id"`
	Name      string `json:"name" firestore:"name"`
	ProjectID string `json:"project_id" firestore:"project_id"`
	Created   *Event `json:"created" firestore:"created"`
	Updated   *Event `json:"updated" firestore:"updated"`
	Deleted   *Event `json:"deleted" firestore:"deleted"`
}

// NewFolderID generates a UUID for folders
func (s *Store) NewFolderID() string {
	return "fol-" + uuid.New().String()
}

// CreateFolder creates a new Folder
func (s *Store) CreateFolder(ctx context.Context, userID string, folder *Folder) error {
	folder.Created = NewEvent(userID)
	_, err := s.Client.Collection(FoldersCollection).Doc(folder.ID).Set(ctx, folder)
	return err
}

// DeleteFolder deletes an existing folder
func (s *Store) DeleteFolder(ctx context.Context, userID string, folder *Folder) error {
	folder.Deleted = NewEvent(userID)
	_, err := s.Client.Collection(FoldersCollection).Doc(folder.ID).Set(ctx, folder)
	return err
}

// UpdateFolder updates an existing folder
func (s *Store) UpdateFolder(ctx context.Context, userID string, folder *Folder) error {
	folder.Updated = NewEvent(userID)
	_, err := s.Client.Collection(FoldersCollection).Doc(folder.ID).Set(ctx, folder)
	return err
}

// GetFolderByID gets a Folder by id
func (s *Store) GetFolderByID(ctx context.Context, id string) (*Folder, error) {
	iter := s.Client.Collection(FoldersCollection).Where("id", "==", id).Limit(1).Documents(ctx)

	snapshot, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	folder := &Folder{}
	snapshot.DataTo(folder)

	if folder.Deleted != nil {
		return nil, nil
	}

	return folder, nil
}
