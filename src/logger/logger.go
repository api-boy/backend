package logger

import (
	"apiboy/backend/src/config"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
)

// Logger wraps the apex logger
type Logger struct {
	Config *config.Config
}

// Field contains a field to include in a log
type Field struct {
	Key string
	Val interface{}
}

// New creates a new logger based on the given level
func New(conf *config.Config) *Logger {
	if conf.UpStage == "development" {
		log.SetHandler(text.Default)
	} else {
		log.SetHandler(json.Default)
	}

	switch conf.LogLevel {
	case "error":
		log.SetLevel(log.ErrorLevel) // error logs
	case "warn":
		log.SetLevel(log.WarnLevel) // warn + error logs
	case "info":
		log.SetLevel(log.InfoLevel) // info + warn + error logs
	case "debug":
		log.SetLevel(log.DebugLevel) // all logs
	default:
		log.SetLevel(log.InfoLevel)
	}

	return &Logger{
		Config: conf,
	}
}

// Error prints an error log
func (l *Logger) Error(msg string, fields ...Field) {
	if len(fields) > 0 {
		log.WithFields(getLogFields(fields)).Error(msg)
	} else {
		log.Error(msg)
	}
}

// Warn prints a warning log
func (l *Logger) Warn(msg string, fields ...Field) {
	if len(fields) > 0 {
		log.WithFields(getLogFields(fields)).Warn(msg)
	} else {
		log.Warn(msg)
	}
}

// Info prints an information log
func (l *Logger) Info(msg string, fields ...Field) {
	if len(fields) > 0 {
		log.WithFields(getLogFields(fields)).Info(msg)
	} else {
		log.Info(msg)
	}
}

// Debug prints a debug log
func (l *Logger) Debug(msg string, fields ...Field) {
	if l.Config.UpStage == "development" {
		if len(fields) > 0 {
			log.WithFields(getLogFields(fields)).Debug(msg)
		} else {
			log.Debug(msg)
		}
	} else if l.Config.LogLevel == "debug" {
		// debug messages are not working properly with the json handler
		// so this will display them as info messages with the [DEBUG] prefix
		msg = "[DEBUG] " + msg

		if len(fields) > 0 {
			log.WithFields(getLogFields(fields)).Info(msg)
		} else {
			log.Info(msg)
		}
	}
}

func getLogFields(fields []Field) log.Fields {
	logFields := log.Fields{}

	for _, f := range fields {
		logFields[f.Key] = f.Val
	}

	return logFields
}
