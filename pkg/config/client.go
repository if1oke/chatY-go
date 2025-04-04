package config

import (
	"github.com/joho/godotenv"
	"log"
)

type IClientConfig interface {
	ServerAddress() string
	ServerPort() string
	ColorScheme() string
	Token() string
	SetServerAddress(serverAddress string)
	SetServerPort(serverPort string)
	SetColorScheme(colorScheme string)
	SetToken(token string)
}

type ClientConfig struct {
	serverAddress string
	serverPort    string
	colorScheme   string
	token         string
}

func (c *ClientConfig) ServerAddress() string {
	return c.serverAddress
}

func (c *ClientConfig) SetServerAddress(serverAddress string) {
	c.serverAddress = serverAddress
}

func (c *ClientConfig) ServerPort() string {
	return c.serverPort
}

func (c *ClientConfig) SetServerPort(serverPort string) {
	c.serverPort = serverPort
}

func (c *ClientConfig) ColorScheme() string {
	return c.colorScheme
}

func (c *ClientConfig) SetColorScheme(colorScheme string) {
	c.colorScheme = colorScheme
}

func (c *ClientConfig) Token() string {
	return c.token
}

func (c *ClientConfig) SetToken(token string) {
	c.token = token
}

func NewClientConfig() *ClientConfig {
	// TODO: сделать сохранение параметров!!!
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	serverAddress := getEnv("SERVER_ADDRESS", "localhost")
	serverPort := getEnv("SERVER_PORT", "8080")
	return &ClientConfig{
		serverAddress: serverAddress,
		serverPort:    serverPort,
	}
}
