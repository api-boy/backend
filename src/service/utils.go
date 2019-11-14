package service

import (
	"context"

	"apiboy/backend/src/enums"
	"apiboy/backend/src/errors"
	"apiboy/backend/src/store"
)

// checkAccessToProject validates if a user has access to a project
func (s *Service) checkAccessToProject(ctx context.Context, userID, projectID string) error {
	// check if a relationship between the project and the user exists
	projectUser, err := s.Store.GetProjectUserByProjectIDAndUserID(ctx, projectID, userID)
	if err != nil {
		return errors.InternalServer{Msg: "Could not get projectUser ", Err: err}
	} else if projectUser == nil {
		return errors.Unauthorized{Msg: "Invalid project for user", Err: err}
	}

	return nil
}

// createExampleProject creates an example project for the given user
func (s *Service) createExampleProject(ctx context.Context, userID string) error {
	// create project
	project := &store.Project{
		ID:   s.Store.NewProjectID(),
		Name: "Example Project",
	}

	if err := s.Store.CreateProject(ctx, userID, project); err != nil {
		return errors.InternalServer{Msg: "Could not create example project", Err: err}
	}

	// create relationship between project and user
	projectUser := &store.ProjectUser{
		ID:        s.Store.NewProjectUserID(project.ID, userID),
		ProjectID: project.ID,
		UserID:    userID,
	}

	if err := s.Store.CreateProjectUser(ctx, userID, projectUser); err != nil {
		return errors.InternalServer{Msg: "Could not create projectUser for example project", Err: err}
	}

	// create folder
	folder := &store.Folder{
		ID:        s.Store.NewFolderID(),
		Name:      "Test Folder",
		ProjectID: project.ID,
	}

	if err := s.Store.CreateFolder(ctx, userID, folder); err != nil {
		return errors.InternalServer{Msg: "Could not create example folder", Err: err}
	}

	// create request
	request := &store.Request{
		ID:        s.Store.NewRequestID(),
		Name:      "Test Request",
		FolderID:  folder.ID,
		ProjectID: project.ID,
		Type:      enums.RequestTypePost,
		URL:       "https://httpbin.org/anything",
		Headers: map[string]string{
			"Accept": "application/json",
		},
		Body: `{
  "key1": "value1",
  "key2": "value2"
}`,
	}

	if err := s.Store.CreateRequest(ctx, userID, request); err != nil {
		return errors.InternalServer{Msg: "Could not create example request", Err: err}
	}

	return nil
}
