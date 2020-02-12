package service

import (
	"apiboy/backend/src/authutils"
	"apiboy/backend/src/errors"
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
)

// SetNewPasswordInput is the input of the endpoint
type SetNewPasswordInput struct {
	Password string `json:"password" validate:"omitempty,min=6"`
	TempCode string `json:"temp_code" validate:"required"`
}

// SetNewPasswordOutput is the output of the endpoint
type SetNewPasswordOutput struct{}

// SetNewPassword implements the business logic for the endpoint
func (s *Service) SetNewPassword(ctx context.Context, input *SetNewPasswordInput) (*SetNewPasswordOutput, error) {
	password := strings.TrimSpace(input.Password)
	tempCode := strings.TrimSpace(input.TempCode)

	decode, err := base64.StdEncoding.DecodeString(tempCode)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not format temp code", Err: err}
	}

	elements := strings.Split(string(decode), "|")
	if len(elements) != 3 {
		return nil, errors.Unauthorized{Msg: "Invalid code"}
	}

	strDateTimeCode := elements[1]
	dateTimeCode, err := time.Parse(time.UnixDate, strDateTimeCode)
	if err != nil { // Always check errors even if they should not happen.
		return nil, errors.InternalServer{Msg: "Could not format date time", Err: err}
	}

	timeNow := time.Now().UTC()
	hrs := timeNow.Sub(dateTimeCode)

	if hrs.Hours() > 24 {
		return nil, errors.Unauthorized{Msg: "Invalid code"}
	}

	userID := elements[0]
	// get user
	user, err := s.Store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not get user", Err: err}
	} else if user == nil {
		return nil, errors.NotFound{Obj: "User"}
	}

	if user.TempCode != tempCode {
		return nil, errors.Unauthorized{Msg: "Invalid code"}
	}

	// hash new password
	hashedPassword, err := authutils.HashPassword(password)
	if err != nil {
		return nil, errors.InternalServer{Msg: "Could not hash password", Err: err}
	}

	user.Password = hashedPassword
	user.TempCode = ""

	if err = s.Store.UpdateUser(ctx, user.ID, user); err != nil {
		return nil, errors.InternalServer{Msg: "Could not update user", Err: err}
	}

	return &SetNewPasswordOutput{}, nil
}

// MakeSetNewPasswordEndpoint creates the endpoint
func MakeSetNewPasswordEndpoint(s *Service, m ...endpoint.Middleware) endpoint.Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		input, ok := request.(*SetNewPasswordInput)
		if !ok {
			return nil, errors.BadRequest{}
		}

		return s.SetNewPassword(ctx, input)
	}

	for _, mw := range m {
		e = mw(e)
	}

	return e
}
