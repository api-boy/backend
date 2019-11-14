package httputils

import (
	"context"
	"net/http"

	"apiboy/backend/src/authutils"

	kithttp "github.com/go-kit/kit/transport/http"
)

type contextKey int

const (
	// ContextKeyRequest is used to store an http request in the context
	ContextKeyRequest contextKey = iota

	// ContextKeyAuthData is used to store the authentication data in the context
	ContextKeyAuthData
)

// SetContextRequest sets the request in the context
func SetContextRequest(ctx context.Context, r *http.Request) context.Context {
	ctx = kithttp.PopulateRequestContext(ctx, r)
	return context.WithValue(ctx, ContextKeyRequest, r)
}

// GetContextRequest gets the request from the context
func GetContextRequest(ctx context.Context) *http.Request {
	return ctx.Value(ContextKeyRequest).(*http.Request)
}

// SetContextAuthData sets the auth data in the context
func SetContextAuthData(ctx context.Context, d *authutils.AuthData) context.Context {
	return context.WithValue(ctx, ContextKeyAuthData, d)
}

// GetContextAuthData gets the auth data from the context
func GetContextAuthData(ctx context.Context) *authutils.AuthData {
	return ctx.Value(ContextKeyAuthData).(*authutils.AuthData)
}
