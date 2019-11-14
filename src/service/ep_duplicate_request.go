package service

import (
	"context"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// DuplicateRequestInput is the input of the endpoint
type DuplicateRequestInput struct {
	ID string `json:"id" validate:"required"`
}

// DuplicateRequestOutput is the output of the endpoint
type DuplicateRequestOutput struct {
	Request *store.Request `json:"request"`
}

// DuplicateRequest implements the business logic for the endpoint
func (s *Service) DuplicateRequest(ctx context.Context, input *DuplicateRequestInput) (*DuplicateRequestOutput, error) {
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

	// duplicate request
	request.ID = s.Store.NewRequestID()
	request.Deleted = nil
	request.Updated = nil
	request.Name += " Copy"

	if err = s.Store.CreateRequest(ctx, authData.UserID, request); err != nil {
		return nil, errors.InternalServer{Msg: "Could not duplicate request", Err: err}
	}

	return &DuplicateRequestOutput{
		Request: request,
	}, nil
}

// MakeDuplicateRequestEndpoint creates the endpoint
func MakeDuplicateRequestEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*DuplicateRequestInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.DuplicateRequest(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
