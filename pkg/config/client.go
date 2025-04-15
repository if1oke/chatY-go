package config

import (
	"github.com/joho/godotenv"
	"log"
)

type ClientConfig struct {
	serverAddress string
	serverPort    string
	username      string
	password      string
}

func (c *ClientConfig) SetUsername(username string) {
	c.username = username
}

func (c *ClientConfig) SetPassword(password string) {
	c.password = password
}

func (c *ClientConfig) Username() string {
	return c.username
}

func (c *ClientConfig) Password() string {
	return c.password
}

func (c *ClientConfig) ServerAddress() string {
	return c.serverAddress
}

func (c *ClientConfig) ServerPort() string {
	return c.serverPort
}

func NewClientConfig() *ClientConfig {
	// TODO: сделать сохранение параметров!!!
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	serverAddress := getEnv("SERVER_ADDRESS", "localhost")
	serverPort := getEnv("SERVER_PORT", "8080")
	username := getEnv("USERNAME", "")
	password := getEnv("PASSWORD", "")
	return &ClientConfig{
		serverAddress: serverAddress,
		serverPort:    serverPort,
		username:      username,
		password:      password,
	}
}
