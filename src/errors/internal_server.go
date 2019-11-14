package errors

import (
	"apiboy/backend/src/logger"
)

// InternalServer indicates an unexpected condition that prevented a service from fulfilling the request
type InternalServer struct {
	Msg string
	Err error
}

// Error returns a string message for this error
func (e InternalServer) Error() string {
	return "Unexpected error"
}

// LogFields returns the fields for logging this error
func (e InternalServer) LogFields() []logger.Field {
	return []logger.Field{
		logger.Field{Key: "Msg", Val: e.Msg},
		logger.Field{Key: "Err", Val: e.Err},
	}
}
