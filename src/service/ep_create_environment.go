package service

import (
	"context"
	"strings"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// CreateEnvironmentInput is the input of the endpoint
type CreateEnvironmentInput struct {
	Name      string            `json:"name" validate:"required"`
	Variables map[string]string `json:"variables" validate:"-"`
	ProjectID string            `json:"project_id" validate:"required"`
}

// CreateEnvironmentOutput is the output of the endpoint
type CreateEnvironmentOutput struct {
	Environment *store.Environment `json:"environment"`
}

// CreateEnvironment implements the business logic for the endpoint
func (s *Service) CreateEnvironment(ctx context.Context, input *CreateEnvironmentInput) (*CreateEnvironmentOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	// check if the user has access to the project
	if err := s.checkAccessToProject(ctx, authData.UserID, input.ProjectID); err != nil {
		return nil, err
	}

	// create environment
	environment := &store.Environment{
		ID:        s.Store.NewEnvironmentID(),
		Name:      strings.TrimSpace(input.Name),
		Variables: input.Variables,
		ProjectID: input.ProjectID,
	}

	if err := s.Store.CreateEnvironment(ctx, authData.UserID, environment); err != nil {
		return nil, errors.InternalServer{Msg: "Could not create environment", Err: err}
	}

	return &CreateEnvironmentOutput{
		Environment: environment,
	}, nil
}

// MakeCreateEnvironmentEndpoint creates the endpoint
func MakeCreateEnvironmentEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*CreateEnvironmentInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.CreateEnvironment(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
