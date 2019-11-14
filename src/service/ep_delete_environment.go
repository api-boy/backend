package service

import (
	"context"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// DeleteEnvironmentInput is the input of the endpoint
type DeleteEnvironmentInput struct {
	ID string `json:"id" validate:"required"`
}

// DeleteEnvironmentOutput is the output of the endpoint
type DeleteEnvironmentOutput struct {
	Environment *store.Environment `json:"environment"`
}

// DeleteEnvironment implements the business logic for the endpoint
func (s *Service) DeleteEnvironment(ctx context.Context, input *DeleteEnvironmentInput) (*DeleteEnvironmentOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	// get environment
	environment, err := s.Store.GetEnvironmentByID(ctx, input.ID)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get environment", Err: err}
	} else if environment == nil {
		return nil, errors.NotFound{Obj: "Environment"}
	}

	// check if the user has access to the project of the environment
	if err := s.checkAccessToProject(ctx, authData.UserID, environment.ProjectID); err != nil {
		return nil, err
	}

	// delete environment
	if err = s.Store.DeleteEnvironment(ctx, authData.UserID, environment); err != nil {
		return nil, errors.InternalServer{Msg: "Could not delete environment", Err: err}
	}

	return &DeleteEnvironmentOutput{
		Environment: environment,
	}, nil
}

// MakeDeleteEnvironmentEndpoint creates the endpoint
func MakeDeleteEnvironmentEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*DeleteEnvironmentInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.DeleteEnvironment(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
