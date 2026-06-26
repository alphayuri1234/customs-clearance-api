package config

import (
	"os"
	"time"
)

type JWTConfig struct {
	SecretKey       string
	ExpirationHours int
}

func LoadJWTConfig() JWTConfig {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "customs-clearance-dev-secret"
	}

	return JWTConfig{
		SecretKey:       secret,
		ExpirationHours: 24,
	}
}

func JWTExpiration() time.Duration {
	return time.Duration(LoadJWTConfig().ExpirationHours) * time.Hour
}
