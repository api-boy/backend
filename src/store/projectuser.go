package store

import (
	"context"

	"google.golang.org/api/iterator"
)

// ProjectUsersCollection is the name of the collection
const ProjectUsersCollection = "projectusers"

// ProjectUser represents a model in the database
type ProjectUser struct {
	ID        string `json:"id" firestore:"id"`
	ProjectID string `json:"project_id" firestore:"project_id"`
	UserID    string `json:"user_id" firestore:"user_id"`
}

// NewProjectUserID generates a UUID for ProjectUser
func (s *Store) NewProjectUserID(projectID, userID string) string {
	return projectID + "-" + userID
}

// CreateProjectUser creates a new ProjectUser
func (s *Store) CreateProjectUser(ctx context.Context, userID string, projectuser *ProjectUser) error {
	_, err := s.Client.Collection(ProjectUsersCollection).Doc(projectuser.ID).Set(ctx, projectuser)
	return err
}

// DeleteProjectUser deletes an existing projectuser
func (s *Store) DeleteProjectUser(ctx context.Context, userID string, projectuser *ProjectUser) error {
	_, err := s.Client.Collection(ProjectUsersCollection).Doc(projectuser.ID).Delete(ctx)
	return err
}

// GetProjectUserByID gets a ProjectUser by id
func (s *Store) GetProjectUserByID(ctx context.Context, id string) (*ProjectUser, error) {
	iter := s.Client.Collection(ProjectUsersCollection).Where("id", "==", id).Limit(1).Documents(ctx)

	snapshot, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	projectuser := &ProjectUser{}
	snapshot.DataTo(projectuser)

	return projectuser, nil
}

// GetProjectUserByProjectIDAndUserID gets a collection ProjectUsers by userid
func (s *Store) GetProjectUserByProjectIDAndUserID(ctx context.Context, projectID string, userID string) (*ProjectUser, error) {
	return s.GetProjectUserByID(ctx, s.NewProjectUserID(projectID, userID))
}
