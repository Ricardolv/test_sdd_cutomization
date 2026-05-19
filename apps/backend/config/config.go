package config

import "os"

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	JWTExpiresIn string
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "9090"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		JWTSecret:    getEnv("JWT_SECRET", "change-me"),
		JWTExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
