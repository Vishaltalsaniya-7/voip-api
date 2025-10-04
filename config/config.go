package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database   DatabaseConfig
	FreeSWITCH FreeSWITCHConfig
	Server     ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type FreeSWITCHConfig struct {
	Host     string
	Port     string
	Password string
}

type ServerConfig struct {
	Port string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "LCXQR7nYUz9PeEXxRYzUlvJous0"),
			DBName:   getEnv("DB_NAME", "fusionpbx"),
		},
		FreeSWITCH: FreeSWITCHConfig{
			Host:     getEnv("FS_HOST", "127.0.0.1"),
			Port:     getEnv("FS_PORT", "8021"),
			Password: getEnv("FS_PASSWORD", "ClueCon"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8086"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
