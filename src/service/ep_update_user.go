package service

import (
	"context"
	"strings"

	"apiboy/backend/src/authutils"
	"apiboy/backend/src/enums"
	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// UpdateUserInput is the input of the endpoint
type UpdateUserInput struct {
	ID       string `json:"id" validate:"-"`
	Name     string `json:"name" validate:"-"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

// UpdateUserOutput is the output of the endpoint
type UpdateUserOutput struct {
	User *store.User `json:"user"`
}

// UpdateUser implements the business logic for the endpoint
func (s *Service) UpdateUser(ctx context.Context, input *UpdateUserInput) (*UpdateUserOutput, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.TrimSpace(input.Email)
	password := strings.TrimSpace(input.Password)

	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	if input.ID == "" {
		input.ID = authData.UserID
	}

	if input.ID != authData.UserID && authData.UserRole != enums.UserRoleAdmin {
		return nil, errors.Unauthorized{}
	}

	// get user
	user, err := s.Store.GetUserByID(ctx, input.ID)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get user", Err: err}
	} else if user == nil {
		return nil, errors.NotFound{Obj: "User"}
	}

	// update user
	if name == "" {
		name = user.Name
	}

	if email == "" {
		email = user.Email
	}

	if password == "" {
		password = user.Password
	} else {
		// hash password
		hashedPassword, err := authutils.HashPassword(password)
		if err != nil {
			return nil, errors.InternalServer{Msg: "Could not hash password", Err: err}
		}

		password = hashedPassword
	}

	user.Name = name
	user.Email = email
	user.Password = password

	if err = s.Store.UpdateUser(ctx, authData.UserID, user); err != nil {
		return nil, errors.InternalServer{Msg: "Could not update user", Err: err}
	}

	return &UpdateUserOutput{
		User: user,
	}, nil
}

// MakeUpdateUserEndpoint creates the endpoint
func MakeUpdateUserEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*UpdateUserInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.UpdateUser(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
