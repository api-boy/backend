package errors

import (
	"fmt"

	"apiboy/backend/src/logger"

	validatorV9 "gopkg.in/go-playground/validator.v9"
)

// InvalidArguments happens when a request contains invalid parameters or is missing required parameters
type InvalidArguments struct {
	Err error
}

// Error returns a string message for this error
func (e InvalidArguments) Error() string {
	msg := "Invalid arguments"

	if ve, ok := e.Err.(validatorV9.ValidationErrors); ok {
		var args []string

		for _, v := range ve {
			args = append(args, v.Field())
		}

		return fmt.Sprintf("%s: %v", msg, args)
	}

	return msg
}

// LogFields returns the fields for logging this error
func (e InvalidArguments) LogFields() []logger.Field {
	return []logger.Field{
		logger.Field{Key: "Err", Val: e.Err},
	}
}
