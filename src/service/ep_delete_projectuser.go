package service

import (
	"context"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"

	"github.com/go-kit/kit/endpoint"
)

// DeleteProjectUserInput is the input of the endpoint
type DeleteProjectUserInput struct {
	ProjectID string `json:"project_id" validate:"required"`
	UserID    string `json:"user_id" validate:"required"`
}

// DeleteProjectUserOutput is the output of the endpoint
type DeleteProjectUserOutput struct{}

// DeleteProjectUser implements the business logic for the endpoint
func (s *Service) DeleteProjectUser(ctx context.Context, input *DeleteProjectUserInput) (*DeleteProjectUserOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	// check if the user has access to the project
	if err := s.checkAccessToProject(ctx, authData.UserID, input.ProjectID); err != nil {
		return nil, err
	}

	// get relationship between project and user
	projectUser, err := s.Store.GetProjectUserByProjectIDAndUserID(ctx, input.ProjectID, input.UserID)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get projectUser ", Err: err}
	} else if projectUser == nil {
		return nil, errors.Unauthorized{Msg: "Invalid project for user", Err: err}
	}

	// delete relationship between project and user
	if err = s.Store.DeleteProjectUser(ctx, authData.UserID, projectUser); err != nil {
		return nil, errors.InternalServer{Msg: "Could not delete projectUser", Err: err}
	}

	return &DeleteProjectUserOutput{}, nil
}

// MakeDeleteProjectUserEndpoint creates the endpoint
func MakeDeleteProjectUserEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*DeleteProjectUserInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.DeleteProjectUser(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
