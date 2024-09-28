package config

import (
	"os"
	"time"
)

type Config struct {
	Level string

	DB     DB
	API    API
	Server Server
}

// type DB represents neccessary data to connect postgres
type DB struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

// type API represents neccessary data for making requests
type API struct {
	Host string
	Port string
}

// type Server represents neccessary data to init server
type Server struct {
	Host        string
	Port        string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

// LoadFromEnv loads config var's from environment
func LoadFromEnv() (*Config, error) {
	cfg := &Config{
		Level: os.Getenv("LOG_LEVEL"),
		DB: DB{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Name: os.Getenv("DB_NAME"),
		},
		API: API{
			Host: os.Getenv("API_HOST"),
			Port: os.Getenv("API_PORT"),
		},
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
	}

	timeout, err := time.ParseDuration(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	idleTimeout, err := time.ParseDuration(os.Getenv("IDLE_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	cfg.Server.Timeout = timeout
	cfg.Server.IdleTimeout = idleTimeout

	return cfg, nil
}
