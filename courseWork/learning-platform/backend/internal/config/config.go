package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret  string
	DBDSN      string
	DockerHost string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		JWTSecret:  getEnv("JWT_SECRET", "fallback_secret_key"),
		DBDSN:      getEnv("DB_DSN", "sqlite.db"),
		DockerHost: getEnv("DOCKER_HOST", "unix:///var/run/docker.sock"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
