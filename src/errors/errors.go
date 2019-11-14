package errors

import (
	"apiboy/backend/src/logger"
)

// Errorer is implemented by all concrete response types that may contain errors
type Errorer interface {
	LogFields() []logger.Field
}
