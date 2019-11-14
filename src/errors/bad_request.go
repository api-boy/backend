package errors

import (
	"fmt"

	"apiboy/backend/src/logger"
)

// BadRequest is returned when a request includes the wrong data
type BadRequest struct {
	Msg string
	Err error
}

// Error returns a string message for this error
func (e BadRequest) Error() string {
	msg := "Invalid request"

	if e.Msg == "" {
		return msg
	}

	return fmt.Sprintf("%s: %v", msg, e.Msg)
}

// LogFields returns the fields for logging this error
func (e BadRequest) LogFields() []logger.Field {
	return []logger.Field{
		logger.Field{Key: "Msg", Val: e.Msg},
		logger.Field{Key: "Err", Val: e.Err},
	}
}
