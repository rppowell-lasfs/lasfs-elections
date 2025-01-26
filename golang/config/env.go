package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Secret    string
	APIConfig APIConfig
	DBConfig  DBConfig
}

type APIConfig struct {
	PublicHost   string
	PublicPort   string
	BcryptSecret string
}

type DBConfig struct {
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
}

// var Envs = initConfig()

var Envs Config

func InitEnvs() {
	Envs = initConfig()
}

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		Secret: getEnv("SECRET", ""),
		APIConfig: APIConfig{
			PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
			PublicPort: getEnv("PUBLIC_PORT", "8080"),
		},
		DBConfig: DBConfig{
			DBUser:     getEnv("DB_USERNAME", "root"),
			DBPassword: getEnv("DB_PASSWORD", "test"),
			DBAddress: fmt.Sprintf("%s:%s",
				getEnv("DB_HOST", "localhost"),
				getEnv("DB_PORT", "3306"),
			),
			DBName: getEnv("DB_NAME", "test"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
