package service

import (
	"context"

	"apiboy/backend/src/enums"
	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// DeleteUserInput is the input of the endpoint
type DeleteUserInput struct {
	ID string `json:"id" validate:"omitempty"`
}

// DeleteUserOutput is the output of the endpoint
type DeleteUserOutput struct {
	User *store.User `json:"user"`
}

// DeleteUser implements the business logic for the endpoint
func (s *Service) DeleteUser(ctx context.Context, input *DeleteUserInput) (*DeleteUserOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	if input.ID == "" {
		input.ID = authData.UserID
	}

	if input.ID != authData.UserID && authData.UserRole != enums.UserRoleAdmin {
		return nil, errors.Unauthorized{}
	}

	// check if the user exists
	user, err := s.Store.GetUserByID(ctx, input.ID)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get user", Err: err}
	} else if user == nil {
		return nil, errors.NotFound{Obj: "User"}
	}

	// delete user
	if err = s.Store.DeleteUser(ctx, authData.UserID, user); err != nil {
		return nil, errors.InternalServer{Msg: "Could not delete user", Err: err}
	}

	return &DeleteUserOutput{
		User: user,
	}, nil
}

// MakeDeleteUserEndpoint creates the endpoint
func MakeDeleteUserEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*DeleteUserInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.DeleteUser(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
