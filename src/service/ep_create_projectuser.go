package service

import (
	"context"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// CreateProjectUserInput is the input of the endpoint
type CreateProjectUserInput struct {
	ProjectID      string `json:"project_id" validate:"required"`
	SharedByUserID string `json:"shared_by_user_id" validate:"required"`
}

// CreateProjectUserOutput is the output of the endpoint
type CreateProjectUserOutput struct {
	ProjectUser *store.ProjectUser `json:"projectuser"`
}

// CreateProjectUser implements the business logic for the endpoint
func (s *Service) CreateProjectUser(ctx context.Context, input *CreateProjectUserInput) (*CreateProjectUserOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	// check if the user who shared the project, has access to it
	if err := s.checkAccessToProject(ctx, input.SharedByUserID, input.ProjectID); err != nil {
		return nil, err
	}

	// check if the relationship between the project and user exists
	projectUser, err := s.Store.GetProjectUserByProjectIDAndUserID(ctx, input.ProjectID, authData.UserID)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get projectUser ", Err: err}
	}

	if projectUser == nil {
		// create relationship between project and user
		projectUser = &store.ProjectUser{
			ID:        s.Store.NewProjectUserID(input.ProjectID, authData.UserID),
			ProjectID: input.ProjectID,
			UserID:    authData.UserID,
		}

		if err = s.Store.CreateProjectUser(ctx, authData.UserID, projectUser); err != nil {
			return nil, errors.InternalServer{Msg: "Could not create projectUser ", Err: err}
		}
	}

	return &CreateProjectUserOutput{
		ProjectUser: projectUser,
	}, nil

}

// MakeCreateProjectUserEndpoint creates the endpoint
func MakeCreateProjectUserEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*CreateProjectUserInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.CreateProjectUser(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
