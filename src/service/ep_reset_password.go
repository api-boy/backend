package service

import (
	"context"
	"strings"

	"apiboy/backend/src/errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

// ResetPasswordInput is the input of the endpoint
type ResetPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordOutput is the output of the endpoint
type ResetPasswordOutput struct{}

// ResetPassword implements the business logic for the endpoint
func (s *Service) ResetPassword(ctx context.Context, input *ResetPasswordInput) (*ResetPasswordOutput, error) {
	email := strings.TrimSpace(input.Email)

	// check if a user with the same email already exists
	user, err := s.Store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get user", Err: err}
	}

	if user == nil {
		return nil, errors.Unauthorized{Msg: "Invalid user"}
	}

	user.TempCode = uuid.New().String()

	if err = s.Store.UpdateUser(ctx, user.ID, user); err != nil {
		return nil, errors.InternalServer{Msg: "Could not generate temp password", Err: err}
	}

	return &ResetPasswordOutput{}, nil
}

// MakeResetPasswordEndpoint creates the endpoint
func MakeResetPasswordEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*ResetPasswordInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.ResetPassword(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
