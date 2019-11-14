package service

import (
	"context"
	"strings"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// CreateProjectInput is the input of the endpoint
type CreateProjectInput struct {
	Name string `json:"name" validate:"required"`
}

// CreateProjectOutput is the output of the endpoint
type CreateProjectOutput struct {
	Project *store.Project `json:"project"`
}

// CreateProject implements the business logic for the endpoint
func (s *Service) CreateProject(ctx context.Context, input *CreateProjectInput) (*CreateProjectOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	// create project
	project := &store.Project{
		ID:   s.Store.NewProjectID(),
		Name: strings.TrimSpace(input.Name),
	}

	if err := s.Store.CreateProject(ctx, authData.UserID, project); err != nil {
		return nil, errors.InternalServer{Msg: "Could not create project", Err: err}
	}

	// create relationship between project and user
	projectUser := &store.ProjectUser{
		ID:        s.Store.NewProjectUserID(project.ID, authData.UserID),
		ProjectID: project.ID,
		UserID:    authData.UserID,
	}

	if err := s.Store.CreateProjectUser(ctx, authData.UserID, projectUser); err != nil {
		return nil, errors.InternalServer{Msg: "Could not create projectUser", Err: err}
	}

	return &CreateProjectOutput{
		Project: project,
	}, nil
}

// MakeCreateProjectEndpoint creates the endpoint
func MakeCreateProjectEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*CreateProjectInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.CreateProject(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
