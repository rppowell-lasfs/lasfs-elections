package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	PublicPort string

	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		PublicPort: getEnv("PUBLIC_PORT", "8080"),

		DBUser:     getEnv("DB_USERNAME", "root"),
		DBPassword: getEnv("DB_PASSWORD", "test"),
		DBAddress: fmt.Sprintf("%s:%s",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "3306"),
		),
		DBName: getEnv("DB_NAME", "test"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
