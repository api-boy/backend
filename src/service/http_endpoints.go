package service

import (
	"github.com/go-kit/kit/endpoint"
)

// HTTPEndpoints collects all of the endpoints that are exposed through http
type HTTPEndpoints struct {
	GetFirebaseCredentialsEndpoint endpoint.Endpoint
	LoginEndpoint                  endpoint.Endpoint
	LogoutEndpoint                 endpoint.Endpoint
	SignupEndpoint                 endpoint.Endpoint
	ResetPasswordEndpoint          endpoint.Endpoint
	SetNewPasswordEndpoint         endpoint.Endpoint
	UpdateUserEndpoint             endpoint.Endpoint
	DeleteUserEndpoint             endpoint.Endpoint
	CreateProjectEndpoint          endpoint.Endpoint
	UpdateProjectEndpoint          endpoint.Endpoint
	DeleteProjectEndpoint          endpoint.Endpoint
	CreateProjectUserEndpoint      endpoint.Endpoint
	DeleteProjectUserEndpoint      endpoint.Endpoint
	CreateFolderEndpoint           endpoint.Endpoint
	DeleteFolderEndpoint           endpoint.Endpoint
	UpdateFolderEndpoint           endpoint.Endpoint
	CreateRequestEndpoint          endpoint.Endpoint
	UpdateRequestEndpoint          endpoint.Endpoint
	DeleteRequestEndpoint          endpoint.Endpoint
	DuplicateRequestEndpoint       endpoint.Endpoint
	CreateEnvironmentEndpoint      endpoint.Endpoint
	UpdateEnvironmentEndpoint      endpoint.Endpoint
	DeleteEnvironmentEndpoint      endpoint.Endpoint
	DuplicateEnvironmentEndpoint   endpoint.Endpoint
}

// MakeHTTPEndpoints returns an HTTPEndpoints struct where each endpoint invokes
// the corresponding method on the provided service
func MakeHTTPEndpoints(s *Service) HTTPEndpoints {
	// Input validation middleware
	vm := s.NewInputValidationMiddleware()

	// Authentication middleware
	am := s.NewAuthMiddleware()

	return HTTPEndpoints{
		GetFirebaseCredentialsEndpoint: MakeGetFirebaseCredentialsEndpoint(s, vm, am),
		LoginEndpoint:                  MakeLoginEndpoint(s, vm),
		LogoutEndpoint:                 MakeLogoutEndpoint(s, vm, am),
		SignupEndpoint:                 MakeSignupEndpoint(s, vm),
		ResetPasswordEndpoint:          MakeResetPasswordEndpoint(s, vm),
		SetNewPasswordEndpoint:         MakeSetNewPasswordEndpoint(s, vm),
		UpdateUserEndpoint:             MakeUpdateUserEndpoint(s, vm, am),
		DeleteUserEndpoint:             MakeDeleteUserEndpoint(s, vm, am),
		CreateProjectEndpoint:          MakeCreateProjectEndpoint(s, vm, am),
		UpdateProjectEndpoint:          MakeUpdateProjectEndpoint(s, vm, am),
		DeleteProjectEndpoint:          MakeDeleteProjectEndpoint(s, vm, am),
		CreateProjectUserEndpoint:      MakeCreateProjectUserEndpoint(s, vm, am),
		DeleteProjectUserEndpoint:      MakeDeleteProjectUserEndpoint(s, vm, am),
		CreateFolderEndpoint:           MakeCreateFolderEndpoint(s, vm, am),
		DeleteFolderEndpoint:           MakeDeleteFolderEndpoint(s, vm, am),
		UpdateFolderEndpoint:           MakeUpdateFolderEndpoint(s, vm, am),
		CreateRequestEndpoint:          MakeCreateRequestEndpoint(s, vm, am),
		UpdateRequestEndpoint:          MakeUpdateRequestEndpoint(s, vm, am),
		DeleteRequestEndpoint:          MakeDeleteRequestEndpoint(s, vm, am),
		DuplicateRequestEndpoint:       MakeDuplicateRequestEndpoint(s, vm, am),
		CreateEnvironmentEndpoint:      MakeCreateEnvironmentEndpoint(s, vm, am),
		UpdateEnvironmentEndpoint:      MakeUpdateEnvironmentEndpoint(s, vm, am),
		DeleteEnvironmentEndpoint:      MakeDeleteEnvironmentEndpoint(s, vm, am),
		DuplicateEnvironmentEndpoint:   MakeDuplicateEnvironmentEndpoint(s, vm, am),
	}
}
