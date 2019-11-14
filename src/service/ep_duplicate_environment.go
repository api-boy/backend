package service

import (
	"context"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// DuplicateEnvironmentInput is the input of the endpoint
type DuplicateEnvironmentInput struct {
	ID string `json:"id" validate:"required"`
}

// DuplicateEnvironmentOutput is the output of the endpoint
type DuplicateEnvironmentOutput struct {
	Environment *store.Environment `json:"environment"`
}

// DuplicateEnvironment implements the business logic for the endpoint
func (s *Service) DuplicateEnvironment(ctx context.Context, input *DuplicateEnvironmentInput) (*DuplicateEnvironmentOutput, error) {
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

	// duplicate environment
	environment.ID = s.Store.NewEnvironmentID()
	environment.Deleted = nil
	environment.Updated = nil
	environment.Name += " Copy"

	if err = s.Store.CreateEnvironment(ctx, authData.UserID, environment); err != nil {
		return nil, errors.InternalServer{Msg: "Could not duplicate environment", Err: err}
	}

	return &DuplicateEnvironmentOutput{
		Environment: environment,
	}, nil
}

// MakeDuplicateEnvironmentEndpoint creates the endpoint
func MakeDuplicateEnvironmentEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*DuplicateEnvironmentInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.DuplicateEnvironment(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
