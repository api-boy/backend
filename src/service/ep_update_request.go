package service

import (
	"context"
	"strings"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// UpdateRequestInput is the input of the endpoint
type UpdateRequestInput struct {
	ID       string            `json:"id" validate:"required"`
	Name     string            `json:"name" validate:"required"`
	FolderID string            `json:"folder_id" validate:"required"`
	Type     string            `json:"type" validate:"omitempty,request_type"`
	URL      string            `json:"url" validate:"-"`
	Headers  map[string]string `json:"headers" validate:"-"`
	Body     string            `json:"body" validate:"-"`
}

// UpdateRequestOutput is the output of the endpoint
type UpdateRequestOutput struct {
	Request *store.Request `json:"request"`
}

// UpdateRequest implements the business logic for the endpoint
func (s *Service) UpdateRequest(ctx context.Context, input *UpdateRequestInput) (*UpdateRequestOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	// get request
	request, err := s.Store.GetRequestByID(ctx, input.ID)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get request", Err: err}
	} else if request == nil {
		return nil, errors.NotFound{Obj: "Request"}
	}

	// check if the user has access to the project of the request
	if err := s.checkAccessToProject(ctx, authData.UserID, request.ProjectID); err != nil {
		return nil, err
	}

	// check if the folder exists and is a folder of the same project (if changed)
	if request.FolderID != input.FolderID {
		folder, err := s.Store.GetFolderByID(ctx, input.FolderID)
		if err != nil {
			return nil, errors.InternalServer{Msg: "Could not get folder", Err: err}
		} else if folder == nil {
			return nil, errors.NotFound{Obj: "Folder"}
		}

		if request.ProjectID != folder.ProjectID {
			return nil, errors.BadRequest{Msg: "Invalid folder for project"}
		}
	}

	// update request
	request.Name = strings.TrimSpace(input.Name)
	request.FolderID = input.FolderID
	request.Type = input.Type
	request.URL = input.URL
	request.Headers = input.Headers
	request.Body = input.Body

	if err = s.Store.UpdateRequest(ctx, authData.UserID, request); err != nil {
		return nil, errors.InternalServer{Msg: "Could not update request", Err: err}
	}

	return &UpdateRequestOutput{
		Request: request,
	}, nil
}

// MakeUpdateRequestEndpoint creates the endpoint
func MakeUpdateRequestEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*UpdateRequestInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.UpdateRequest(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
