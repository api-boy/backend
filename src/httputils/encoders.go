package httputils

import (
	"context"
	"encoding/json"
	"net/http"

	"apiboy/backend/src/errors"
	"apiboy/backend/src/logger"

	kithttp "github.com/go-kit/kit/transport/http"
)

// ResponseEncoder is a standard encoder for JSON responses
func ResponseEncoder(log *logger.Logger) kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		return json.NewEncoder(w).Encode(response)
	}
}

// ErrorEncoder sends a formatted error back to the ResponseWriter
func ErrorEncoder(log *logger.Logger) kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		msg := err.Error()

		if e, ok := err.(errors.Errorer); ok {
			fields := e.LogFields()
			log.Debug(msg, fields...)
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(statusCodeForError(err))

		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": msg,
		})
	}
}

// statusCodeForError returns the http Status code for a given error
func statusCodeForError(err error) int {
	switch err.(type) {
	case errors.NotFound:
		return http.StatusNotFound
	case errors.Unauthenticated:
		return http.StatusUnauthorized
	case errors.Unauthorized:
		return http.StatusForbidden
	case errors.InvalidArguments:
		return http.StatusBadRequest
	case errors.BadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
