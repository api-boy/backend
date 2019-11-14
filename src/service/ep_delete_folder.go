package service

import (
	"context"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// DeleteFolderInput is the input of the endpoint
type DeleteFolderInput struct {
	ID string `json:"id" validate:"required"`
}

// DeleteFolderOutput is the output of the endpoint
type DeleteFolderOutput struct {
	Folder *store.Folder `json:"folder"`
}

// DeleteFolder implements the business logic for the endpoint
func (s *Service) DeleteFolder(ctx context.Context, input *DeleteFolderInput) (*DeleteFolderOutput, error) {
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

	// delete folder
	if err = s.Store.DeleteFolder(ctx, authData.UserID, folder); err != nil {
		return nil, errors.InternalServer{Msg: "Could not delete folder", Err: err}
	}

	return &DeleteFolderOutput{
		Folder: folder,
	}, nil
}

// MakeDeleteFolderEndpoint creates the endpoint
func MakeDeleteFolderEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*DeleteFolderInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.DeleteFolder(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
