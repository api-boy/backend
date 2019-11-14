package httputils

import (
	"context"
	"net/http"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/logger"
)

type metadataOutput struct {
	Path     string `json:"path"`
	Method   string `json:"method"`
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
}

// MetadataHandler returns an http handler for metadata requests
func MetadataHandler(ctx context.Context, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := metadataOutput{
			Path:     r.URL.Path,
			Method:   r.Method,
			Protocol: r.Proto,
			Host:     r.Host,
		}

		ResponseEncoder(log)(ctx, w, res)
	}
}

// NotFoundHandler returns an http handler for not found requests
func NotFoundHandler(ctx context.Context, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := errors.NotFound{
			Obj: r.URL.Path,
		}

		ErrorEncoder(log)(ctx, err, w)
	}
}

// CORSHandler returns an http handler for CORS requests
func CORSHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	}
}
