package errors

import (
	"fmt"

	"apiboy/backend/src/logger"
)

// NotFound is returned when a requested object is not found
type NotFound struct {
	Obj string
	Err error
}

// Error returns a string message for this error
func (e NotFound) Error() string {
	msg := "Resource not found"

	if e.Obj == "" {
		return msg
	}

	return fmt.Sprintf("%s: %v", msg, e.Obj)
}

// LogFields returns the fields for logging this error
func (e NotFound) LogFields() []logger.Field {
	return []logger.Field{
		logger.Field{Key: "Obj", Val: e.Obj},
		logger.Field{Key: "Err", Val: e.Err},
	}
}
