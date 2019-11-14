package service

import (
	"context"
	"strings"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// UpdateEnvironmentInput is the input of the endpoint
type UpdateEnvironmentInput struct {
	ID        string            `json:"id" validate:"required"`
	Name      string            `json:"name" validate:"required"`
	Variables map[string]string `json:"variables" validate:"-"`
}

// UpdateEnvironmentOutput is the output of the endpoint
type UpdateEnvironmentOutput struct {
	Environment *store.Environment `json:"environment"`
}

// UpdateEnvironment implements the business logic for the endpoint
func (s *Service) UpdateEnvironment(ctx context.Context, input *UpdateEnvironmentInput) (*UpdateEnvironmentOutput, error) {
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

	// update environment
	environment.Name = strings.TrimSpace(input.Name)
	environment.Variables = input.Variables

	if err = s.Store.UpdateEnvironment(ctx, authData.UserID, environment); err != nil {
		return nil, errors.InternalServer{Msg: "Could not update environment", Err: err}
	}

	return &UpdateEnvironmentOutput{
		Environment: environment,
	}, nil
}

// MakeUpdateEnvironmentEndpoint creates the endpoint
func MakeUpdateEnvironmentEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*UpdateEnvironmentInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.UpdateEnvironment(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
