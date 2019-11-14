package service

import (
	"context"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// DeleteProjectInput is the input of the endpoint
type DeleteProjectInput struct {
	ID string `json:"id" validate:"required"`
}

// DeleteProjectOutput is the output of the endpoint
type DeleteProjectOutput struct {
	Project *store.Project `json:"project"`
}

// DeleteProject implements the business logic for the endpoint
func (s *Service) DeleteProject(ctx context.Context, input *DeleteProjectInput) (*DeleteProjectOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	// get project
	project, err := s.Store.GetProjectByID(ctx, input.ID)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get project", Err: err}
	} else if project == nil {
		return nil, errors.NotFound{Obj: "Project"}
	}

	// check if the user is the owner of the project
	if project.Created.By != authData.UserID {
		return nil, errors.Unauthorized{Msg: "The user is not the owner of the project", Err: err}
	}

	// delete project
	if err = s.Store.DeleteProject(ctx, authData.UserID, project); err != nil {
		return nil, errors.InternalServer{Msg: "Could not delete project", Err: err}
	}

	return &DeleteProjectOutput{
		Project: project,
	}, nil
}

// MakeDeleteProjectEndpoint creates the endpoint
func MakeDeleteProjectEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*DeleteProjectInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.DeleteProject(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
