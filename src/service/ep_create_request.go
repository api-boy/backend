package service

import (
	"context"
	"strings"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// CreateRequestInput is the input of the endpoint
type CreateRequestInput struct {
	Name     string            `json:"name" validate:"required"`
	FolderID string            `json:"folder_id" validate:"required"`
	Type     string            `json:"type" validate:"omitempty,request_type"`
	URL      string            `json:"url" validate:"-"`
	Headers  map[string]string `json:"headers" validate:"-"`
	Body     string            `json:"body" validate:"-"`
}

// CreateRequestOutput is the output of the endpoint
type CreateRequestOutput struct {
	Request *store.Request `json:"request"`
}

// CreateRequest implements the business logic for the endpoint
func (s *Service) CreateRequest(ctx context.Context, input *CreateRequestInput) (*CreateRequestOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	// get folder
	folder, err := s.Store.GetFolderByID(ctx, input.FolderID)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get folder", Err: err}
	} else if folder == nil {
		return nil, errors.NotFound{Obj: "Folder"}
	}

	// check if the user has access to the project of the folder
	if err := s.checkAccessToProject(ctx, authData.UserID, folder.ProjectID); err != nil {
		return nil, err
	}

	// create request
	request := &store.Request{
		ID:        s.Store.NewRequestID(),
		Name:      strings.TrimSpace(input.Name),
		FolderID:  input.FolderID,
		ProjectID: folder.ProjectID,
		Type:      input.Type,
		URL:       input.URL,
		Headers:   input.Headers,
		Body:      input.Body,
	}

	if err = s.Store.CreateRequest(ctx, authData.UserID, request); err != nil {
		return nil, errors.InternalServer{Msg: "Could not create request", Err: err}
	}

	return &CreateRequestOutput{
		Request: request,
	}, nil
}

// MakeCreateRequestEndpoint creates the endpoint
func MakeCreateRequestEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*CreateRequestInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.CreateRequest(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
