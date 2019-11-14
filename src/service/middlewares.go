package service

import (
	"context"
	"strings"
	"time"

	"apiboy/backend/src/authutils"
	"apiboy/backend/src/enums"
	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"

	"github.com/go-kit/kit/endpoint"
	validatorV9 "gopkg.in/go-playground/validator.v9"
)

// NewAuthMiddleware returns an endpoint middleware to handle authentication
func (s *Service) NewAuthMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// get http request from context
			r := httputils.GetContextRequest(ctx)

			// get token string from the auth header
			authHeader := r.Header.Get("Authorization")
			jwt := strings.Replace(authHeader, "Bearer ", "", 1)

			// parse jwt
			authData, err := authutils.ParseJWT(s.Config.JWTSignKey, jwt)
			if err != nil {
				return nil, errors.Unauthenticated{Msg: "Could not parse jwt", Err: err}
			}

			// check if the token is valid for the user in the database
			token, err := s.Store.GetTokenByID(ctx, authData.JwtID)
			if err != nil || token == nil {
				return nil, errors.Unauthenticated{Msg: "Could not get token", Err: err}
			}

			if token.UserID != authData.UserID {
				return nil, errors.Unauthenticated{Msg: "Invalid token for the user"}
			}

			// populate context with auth data
			ctx = httputils.SetContextAuthData(ctx, authData)
			return next(ctx, request)
		}
	}
}

// NewInputValidationMiddleware returns an endpoint middleware to handle input validations
func (s *Service) NewInputValidationMiddleware() endpoint.Middleware {
	inputValidator := validatorV9.New()

	// register custom validations
	inputValidator.RegisterValidation("time_rfc3339", func(fl validatorV9.FieldLevel) bool {
		value := fl.Field().String()

		if _, err := time.Parse(time.RFC3339, value); err != nil {
			return false
		}

		return true
	})

	inputValidator.RegisterValidation("request_type", func(fl validatorV9.FieldLevel) bool {
		value := fl.Field().String()

		return enums.IsValidRequestType(value)
	})

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err := inputValidator.Struct(request); err != nil {
				return nil, errors.InvalidArguments{Err: err}
			}

			return next(ctx, request)
		}
	}
}
