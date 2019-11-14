package service

import (
	"context"
	"strings"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// UpdateProjectInput is the input of the endpoint
type UpdateProjectInput struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

// UpdateProjectOutput is the output of the endpoint
type UpdateProjectOutput struct {
	Project *store.Project `json:"project"`
}

// UpdateProject implements the business logic for the endpoint
func (s *Service) UpdateProject(ctx context.Context, input *UpdateProjectInput) (*UpdateProjectOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	// check if the user has access to the project
	if err := s.checkAccessToProject(ctx, authData.UserID, input.ID); err != nil {
		return nil, err
	}

	// get project
	project, err := s.Store.GetProjectByID(ctx, input.ID)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get project", Err: err}
	} else if project == nil {
		return nil, errors.NotFound{Obj: "Project"}
	}

	// update project
	project.Name = strings.TrimSpace(input.Name)

	if err = s.Store.UpdateProject(ctx, authData.UserID, project); err != nil {
		return nil, errors.InternalServer{Msg: "Could not update project", Err: err}
	}

	return &UpdateProjectOutput{
		Project: project,
	}, nil
}

// MakeUpdateProjectEndpoint creates the endpoint
func MakeUpdateProjectEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*UpdateProjectInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.UpdateProject(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
