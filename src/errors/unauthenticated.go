package errors

import (
	"apiboy/backend/src/logger"
)

// Unauthenticated happens when a request is not authorized
type Unauthenticated struct {
	Msg string
	Err error
}

// Error returns a string message for this error
func (e Unauthenticated) Error() string {
	return "Unauthenticated request"
}

// LogFields returns the fields for logging this error
func (e Unauthenticated) LogFields() []logger.Field {
	return []logger.Field{
		logger.Field{Key: "Msg", Val: e.Msg},
		logger.Field{Key: "Err", Val: e.Err},
	}
}
