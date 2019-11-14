package service

import (
	"context"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// DeleteRequestInput is the input of the endpoint
type DeleteRequestInput struct {
	ID string `json:"id" validate:"required"`
}

// DeleteRequestOutput is the output of the endpoint
type DeleteRequestOutput struct {
	Request *store.Request `json:"request"`
}

// DeleteRequest implements the business logic for the endpoint
func (s *Service) DeleteRequest(ctx context.Context, input *DeleteRequestInput) (*DeleteRequestOutput, error) {
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

	// delete request
	if err = s.Store.DeleteRequest(ctx, authData.UserID, request); err != nil {
		return nil, errors.InternalServer{Msg: "Could not delete request", Err: err}
	}

	return &DeleteRequestOutput{
		Request: request,
	}, nil
}

// MakeDeleteRequestEndpoint creates the endpoint
func MakeDeleteRequestEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*DeleteRequestInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.DeleteRequest(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
