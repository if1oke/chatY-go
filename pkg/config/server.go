package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type IConfig interface {
	ServerPort() string
	ServerAddress() string
	AuthPort() string
}

type Config struct {
	serverPort    string
	serverAddress string
	authPort      string
}

func (c Config) ServerPort() string {
	return c.serverPort
}

func (c Config) ServerAddress() string {
	return c.serverAddress
}

func (c Config) AuthPort() string {
	return c.authPort
}

func NewConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	serverPort := getEnv("SERVER_PORT", "33100")
	serverAddress := getEnv("SERVER_ADDRESS", "localhost")
	authPort := getEnv("AUTH_PORT", "33101")
	return &Config{serverPort, serverAddress, authPort}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
