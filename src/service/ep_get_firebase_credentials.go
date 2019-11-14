package service

import (
	"context"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/httputils"

	"firebase.google.com/go/auth"
	"github.com/go-kit/kit/endpoint"
)

// GetFirebaseCredentialsInput is the input of the endpoint
type GetFirebaseCredentialsInput struct{}

// GetFirebaseCredentialsOutput is the output of the endpoint
type GetFirebaseCredentialsOutput struct {
	ProjectID   string `json:"project_id"`
	APIKey      string `json:"api_key"`
	AccessToken string `json:"access_token"`
}

// GetFirebaseCredentials implements the business logic for the endpoint
func (s *Service) GetFirebaseCredentials(ctx context.Context, input *GetFirebaseCredentialsInput) (*GetFirebaseCredentialsOutput, error) {
	// get the auth data from the context
	authData := httputils.GetContextAuthData(ctx)

	// create the user in firebase auth if not exists
	usr, err := s.FirebaseAuthClient.GetUser(ctx, authData.UserID)
	if err == nil {
		if usr.Email != authData.UserEmail {
			// update user in firebase auth
			user := (&auth.UserToUpdate{}).Email(authData.UserEmail)

			_, err := s.FirebaseAuthClient.UpdateUser(ctx, authData.UserID, user)
			if err != nil {
				return nil, errors.InternalServer{Msg: "Could not update user in firebase auth", Err: err}
			}
		}
	} else if auth.IsUserNotFound(err) {
		// create user in firebase auth
		user := (&auth.UserToCreate{}).UID(authData.UserID).Email(authData.UserEmail)

		_, err := s.FirebaseAuthClient.CreateUser(ctx, user)
		if err != nil {
			return nil, errors.InternalServer{Msg: "Could not create user in firebase auth", Err: err}
		}
	} else if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get user from firebase auth", Err: err}
	}

	// create firebase auth token
	claims := map[string]interface{}{
		"user_name":  authData.UserName,
		"user_email": authData.UserEmail,
		"user_role":  authData.UserRole,
	}

	accessToken, err := s.FirebaseAuthClient.CustomTokenWithClaims(ctx, authData.UserID, claims)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not create firebase auth token", Err: err}
	}

	return &GetFirebaseCredentialsOutput{
		ProjectID:   s.Config.FirebaseProjectID,
		APIKey:      s.Config.FirebaseAPIKey,
		AccessToken: accessToken,
	}, nil
}

// MakeGetFirebaseCredentialsEndpoint creates the endpoint
func MakeGetFirebaseCredentialsEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*GetFirebaseCredentialsInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.GetFirebaseCredentials(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
