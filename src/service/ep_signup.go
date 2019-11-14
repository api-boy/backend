package service

import (
	"context"
	"strings"

	"apiboy/backend/src/authutils"
	"apiboy/backend/src/enums"
	"apiboy/backend/src/errors"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// SignupInput is the input of the endpoint
type SignupInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// SignupOutput is the output of the endpoint
type SignupOutput struct {
	JWT string `json:"jwt"`
}

// Signup implements the business logic for the endpoint
func (s *Service) Signup(ctx context.Context, input *SignupInput) (*SignupOutput, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.TrimSpace(input.Email)
	password := strings.TrimSpace(input.Password)

	// check if a user with the same email already exists
	user, err := s.Store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get user", Err: err}
	} else if user != nil {
		return nil, errors.BadRequest{Msg: "User already exists"}
	}

	// hash password
	hashedPassword, err := authutils.HashPassword(password)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not hash password", Err: err}
	}

	// create user
	user = &store.User{
		ID:       s.Store.NewUserID(),
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     enums.UserRoleUser,
	}

	if err = s.Store.CreateUser(ctx, user.ID, user); err != nil {
		return nil, errors.InternalServer{Msg: "Could not create user", Err: err}
	}

	// create example project for the new user
	if err = s.createExampleProject(ctx, user.ID); err != nil {
		return nil, err
	}

	// create token for the user
	token := &store.Token{
		ID:     s.Store.NewTokenID(),
		UserID: user.ID,
	}

	if err = s.Store.CreateToken(ctx, token); err != nil {
		return nil, errors.InternalServer{Msg: "Could not create token", Err: err}
	}

	// create authentication jwt
	authData := &authutils.AuthData{
		JwtID:     token.ID,
		UserID:    user.ID,
		UserName:  user.Name,
		UserEmail: user.Email,
		UserRole:  user.Role,
	}

	jwt, err := authutils.NewJWT(s.Config.JWTIssuer, s.Config.JWTSignKey, authData)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not create jwt", Err: err}
	}

	return &SignupOutput{
		JWT: jwt,
	}, nil
}

// MakeSignupEndpoint creates the endpoint
func MakeSignupEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*SignupInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.Signup(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
