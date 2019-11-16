package service

import (
	"context"
	"net/http"

	"apiboy/backend/src/httputils"
	"apiboy/backend/src/logger"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler
func MakeHTTPHandler(ctx context.Context, log *logger.Logger, e HTTPEndpoints) http.Handler {
	r := mux.NewRouter()

	defaultOptions := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(httputils.ErrorEncoder(log)),
		kithttp.ServerBefore(httputils.SetContextRequest),
	}

	// CORS Handler
	r.Methods("OPTIONS").HandlerFunc(httputils.CORSHandler()).Name("CORS")

	// Metadata Handler
	r.Path("/").HandlerFunc(httputils.MetadataHandler(ctx, log)).Name("Metadata")

	/*******************************************/
	/*** Setup HTTP routes for the endpoints ***/
	/*******************************************/

	r.Methods("POST").Path("/auth/get_firebase_credentials").Handler(kithttp.NewServer(
		e.GetFirebaseCredentialsEndpoint,
		httputils.DecodeRPCRequest(&GetFirebaseCredentialsInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("GetFirebaseCredentials")

	r.Methods("POST").Path("/auth/login").Handler(kithttp.NewServer(
		e.LoginEndpoint,
		httputils.DecodeRPCRequest(&LoginInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("Login")

	r.Methods("POST").Path("/auth/logout").Handler(kithttp.NewServer(
		e.LogoutEndpoint,
		httputils.DecodeRPCRequest(&LogoutInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("Logout")

	r.Methods("POST").Path("/auth/signup").Handler(kithttp.NewServer(
		e.SignupEndpoint,
		httputils.DecodeRPCRequest(&SignupInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("Signup")

	r.Methods("POST").Path("/auth/reset_password").Handler(kithttp.NewServer(
		e.ResetPasswordEndpoint,
		httputils.DecodeRPCRequest(&ResetPasswordInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("ResetPassword")

	r.Methods("POST").Path("/users/update").Handler(kithttp.NewServer(
		e.UpdateUserEndpoint,
		httputils.DecodeRPCRequest(&UpdateUserInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("UpdateUser")

	r.Methods("POST").Path("/users/delete").Handler(kithttp.NewServer(
		e.DeleteUserEndpoint,
		httputils.DecodeRPCRequest(&DeleteUserInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("DeleteUser")

	r.Methods("POST").Path("/projects/create").Handler(kithttp.NewServer(
		e.CreateProjectEndpoint,
		httputils.DecodeRPCRequest(&CreateProjectInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("CreateProject")

	r.Methods("POST").Path("/projects/update").Handler(kithttp.NewServer(
		e.UpdateProjectEndpoint,
		httputils.DecodeRPCRequest(&UpdateProjectInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("UpdateProject")

	r.Methods("POST").Path("/projects/delete").Handler(kithttp.NewServer(
		e.DeleteProjectEndpoint,
		httputils.DecodeRPCRequest(&DeleteProjectInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("DeleteProject")

	r.Methods("POST").Path("/projects-users/create").Handler(kithttp.NewServer(
		e.CreateProjectUserEndpoint,
		httputils.DecodeRPCRequest(&CreateProjectUserInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("CreateProjectUser")

	r.Methods("POST").Path("/projects-users/delete").Handler(kithttp.NewServer(
		e.DeleteProjectUserEndpoint,
		httputils.DecodeRPCRequest(&DeleteProjectUserInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("DeleteProjectUser")

	r.Methods("POST").Path("/folders/create").Handler(kithttp.NewServer(
		e.CreateFolderEndpoint,
		httputils.DecodeRPCRequest(&CreateFolderInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("CreateFolder")

	r.Methods("POST").Path("/folders/delete").Handler(kithttp.NewServer(
		e.DeleteFolderEndpoint,
		httputils.DecodeRPCRequest(&DeleteFolderInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("DeleteFolder")

	r.Methods("POST").Path("/folders/update").Handler(kithttp.NewServer(
		e.UpdateFolderEndpoint,
		httputils.DecodeRPCRequest(&UpdateFolderInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("UpdateFolder")

	r.Methods("POST").Path("/requests/create").Handler(kithttp.NewServer(
		e.CreateRequestEndpoint,
		httputils.DecodeRPCRequest(&CreateRequestInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("CreateRequest")

	r.Methods("POST").Path("/requests/update").Handler(kithttp.NewServer(
		e.UpdateRequestEndpoint,
		httputils.DecodeRPCRequest(&UpdateRequestInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("UpdateRequest")

	r.Methods("POST").Path("/requests/delete").Handler(kithttp.NewServer(
		e.DeleteRequestEndpoint,
		httputils.DecodeRPCRequest(&DeleteRequestInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("DeleteRequest")

	r.Methods("POST").Path("/requests/duplicate").Handler(kithttp.NewServer(
		e.DuplicateRequestEndpoint,
		httputils.DecodeRPCRequest(&DuplicateRequestInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("DuplicateRequest")

	r.Methods("POST").Path("/environments/create").Handler(kithttp.NewServer(
		e.CreateEnvironmentEndpoint,
		httputils.DecodeRPCRequest(&CreateEnvironmentInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("CreateEnvironment")

	r.Methods("POST").Path("/environments/update").Handler(kithttp.NewServer(
		e.UpdateEnvironmentEndpoint,
		httputils.DecodeRPCRequest(&UpdateEnvironmentInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("UpdateEnvironment")

	r.Methods("POST").Path("/environments/delete").Handler(kithttp.NewServer(
		e.DeleteEnvironmentEndpoint,
		httputils.DecodeRPCRequest(&DeleteEnvironmentInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("DeleteEnvironment")

	r.Methods("POST").Path("/environments/duplicate").Handler(kithttp.NewServer(
		e.DuplicateEnvironmentEndpoint,
		httputils.DecodeRPCRequest(&DuplicateEnvironmentInput{}),
		httputils.ResponseEncoder(log),
		defaultOptions...,
	)).Name("DuplicateEnvironment")

	/*******************************************/

	// NotFound Handler: catch any other request with this handler
	r.PathPrefix("/").HandlerFunc(httputils.NotFoundHandler(ctx, log)).Name("NotFound")

	return r
}
