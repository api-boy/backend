package service

import (
	"context"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"

	"github.com/go-kit/kit/endpoint"
)

// LogoutInput is the input of the endpoint
type LogoutInput struct{}

// LogoutOutput is the output of the endpoint
type LogoutOutput struct{}

// Logout implements the business logic for the endpoint
func (s *Service) Logout(ctx context.Context, input *LogoutInput) (*LogoutOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	// delete the token
	if err := s.Store.DeleteToken(ctx, authData.JwtID); err != nil {
		return nil, errors.InternalServer{Msg: "Could not delete token", Err: err}
	}

	return &LogoutOutput{}, nil
}

// MakeLogoutEndpoint creates the endpoint
func MakeLogoutEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*LogoutInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.Logout(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
