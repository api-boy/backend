package config

import (
	"os"
)

// Config contains the configuration parameters for the app
type Config struct {
	UpStage           string
	Port              string
	LogLevel          string
	FirebaseProjectID string
	FirebaseAPIKey    string
	JWTIssuer         string
	JWTSignKey        string
}

// New reads the app configurationa
func New() *Config {
	return &Config{
		UpStage:           os.Getenv("UP_STAGE"),
		Port:              os.Getenv("PORT"),
		LogLevel:          os.Getenv("LOG_LEVEL"),
		FirebaseProjectID: os.Getenv("FIREBASE_PROJECT_ID"),
		FirebaseAPIKey:    os.Getenv("FIREBASE_API_KEY"),
		JWTIssuer:         os.Getenv("JWT_ISSUER"),
		JWTSignKey:        os.Getenv("JWT_SIGN_KEY"),
	}
}
