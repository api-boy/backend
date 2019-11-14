package store

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// ProjectsCollection is the name of the collection
const ProjectsCollection = "projects"

// Project represents a model in the database
type Project struct {
	ID      string `json:"id" firestore:"id"`
	Name    string `json:"name" firestore:"name"`
	Created *Event `json:"created" firestore:"created"`
	Updated *Event `json:"updated" firestore:"updated"`
	Deleted *Event `json:"deleted" firestore:"deleted"`
}

// NewProjectID generates a UUID for Projects
func (s *Store) NewProjectID() string {
	return "pro-" + uuid.New().String()
}

// CreateProject creates a new Project
func (s *Store) CreateProject(ctx context.Context, userID string, project *Project) error {
	project.Created = NewEvent(userID)
	_, err := s.Client.Collection(ProjectsCollection).Doc(project.ID).Set(ctx, project)
	return err
}

// UpdateProject updates an existing Project
func (s *Store) UpdateProject(ctx context.Context, userID string, project *Project) error {
	project.Updated = NewEvent(userID)
	_, err := s.Client.Collection(ProjectsCollection).Doc(project.ID).Set(ctx, project)
	return err
}

// DeleteProject deletes an existing project
func (s *Store) DeleteProject(ctx context.Context, userID string, project *Project) error {
	project.Deleted = NewEvent(userID)
	_, err := s.Client.Collection(ProjectsCollection).Doc(project.ID).Set(ctx, project)
	return err
}

// GetProjectByID gets a project by id
func (s *Store) GetProjectByID(ctx context.Context, id string) (*Project, error) {
	iter := s.Client.Collection(ProjectsCollection).Where("id", "==", id).Limit(1).Documents(ctx)

	snapshot, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	project := &Project{}
	snapshot.DataTo(project)

	if project.Deleted != nil {
		return nil, nil
	}

	return project, nil
}
