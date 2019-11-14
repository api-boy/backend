package service

import (
	"context"
	"strings"

	"apiboy/backend/src/authutils"
	"apiboy/backend/src/errors"
	"apiboy/backend/src/store"

	"github.com/go-kit/kit/endpoint"
)

// LoginInput is the input of the endpoint
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginOutput is the output of the endpoint
type LoginOutput struct {
	JWT string `json:"jwt"`
}

// Login implements the business logic for the endpoint
func (s *Service) Login(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	email := strings.TrimSpace(input.Email)
	password := strings.TrimSpace(input.Password)

	// get the user with the given email
	user, err := s.Store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get user", Err: err}
	}

	if user == nil {
		return nil, errors.Unauthorized{Msg: "Invalid user"}
	}

	// check user password
	if err = authutils.CheckPassword(user.Password, password); err != nil {
		return nil, errors.Unauthorized{Msg: "Invalid password", Err: err}
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

	return &LoginOutput{
		JWT: jwt,
	}, nil
}

// MakeLoginEndpoint creates the endpoint
func MakeLoginEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*LoginInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.Login(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
