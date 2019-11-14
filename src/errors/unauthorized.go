package errors

import (
	"apiboy/backend/src/logger"
)

// Unauthorized indicates an unauthorized access
type Unauthorized struct {
	Msg string
	Err error
}

// Error returns a string message for this error
func (e Unauthorized) Error() string {
	return "Unauthorized action"
}

// LogFields returns the fields for logging this error
func (e Unauthorized) LogFields() []logger.Field {
	return []logger.Field{
		logger.Field{Key: "Msg", Val: e.Msg},
		logger.Field{Key: "Err", Val: e.Err},
	}
}
