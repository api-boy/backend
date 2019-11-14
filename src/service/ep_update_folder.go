package service

import (
	"context"
	"strings"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// UpdateFolderInput is the input of the endpoint
type UpdateFolderInput struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

// UpdateFolderOutput is the output of the endpoint
type UpdateFolderOutput struct {
	Folder *store.Folder `json:"folder"`
}

// UpdateFolder implements the business logic for the endpoint
func (s *Service) UpdateFolder(ctx context.Context, input *UpdateFolderInput) (*UpdateFolderOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	// get folder
	folder, err := s.Store.GetFolderByID(ctx, input.ID)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get folder", Err: err}
	} else if folder == nil {
		return nil, errors.NotFound{Obj: "Folder"}
	}

	// check if the user has access to the project of the folder
	if err := s.checkAccessToProject(ctx, authData.UserID, folder.ProjectID); err != nil {
		return nil, err
	}

	// update folder
	folder.Name = strings.TrimSpace(input.Name)

	if err = s.Store.UpdateFolder(ctx, authData.UserID, folder); err != nil {
		return nil, errors.InternalServer{Msg: "Could not update folder", Err: err}
	}

	return &UpdateFolderOutput{
		Folder: folder,
	}, nil
}

// MakeUpdateFolderEndpoint creates the endpoint
func MakeUpdateFolderEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*UpdateFolderInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.UpdateFolder(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
