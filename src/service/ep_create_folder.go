package service

import (
	"context"
	"strings"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// CreateFolderInput is the input of the endpoint
type CreateFolderInput struct {
	Name      string `json:"name" validate:"required"`
	ProjectID string `json:"project_id" validate:"required"`
}

// CreateFolderOutput is the output of the endpoint
type CreateFolderOutput struct {
	Folder *store.Folder `json:"folder"`
}

// CreateFolder implements the business logic for the endpoint
func (s *Service) CreateFolder(ctx context.Context, input *CreateFolderInput) (*CreateFolderOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	// check if the user has access to the project
	if err := s.checkAccessToProject(ctx, authData.UserID, input.ProjectID); err != nil {
		return nil, err
	}

	// create folder
	folder := &store.Folder{
		ID:        s.Store.NewFolderID(),
		Name:      strings.TrimSpace(input.Name),
		ProjectID: input.ProjectID,
	}

	if err := s.Store.CreateFolder(ctx, authData.UserID, folder); err != nil {
		return nil, errors.InternalServer{Msg: "Could not create folder", Err: err}
	}

	return &CreateFolderOutput{
		Folder: folder,
	}, nil
}

// MakeCreateFolderEndpoint creates the endpoint
func MakeCreateFolderEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*CreateFolderInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.CreateFolder(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
