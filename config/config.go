package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	GRPCPort          string
	GRPCServerAddress string
	HTTPServerAddress string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		log.Println("Continuing with existing environment variables")
	}

	return Config{
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		GRPCPort:          os.Getenv("GRPC_PORT"),
		GRPCServerAddress: os.Getenv("GRPC_SERVER_ADDRESS"),
		HTTPServerAddress: os.Getenv("HTTP_SERVER_ADDRESS"),
	}
}
