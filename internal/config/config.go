package config

import (
	"log/slog"
	"os"
	"time"
)

type Config struct {
	Level      string
	UseTestApi bool

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

// AsLogValue represents DB struct as slog.Value
// Used for logging
func (db *DB) AsLogValue() slog.Value {
	return slog.GroupValue(
		slog.String("host", db.Host),
		slog.String("port", db.Port),
		slog.String("user", db.User),
		slog.String("pass", "[hidden]"),
		slog.String("name", db.Name),
	)
}

// AsLogValue represents API struct as slog.Value
// Used for logging
func (a *API) AsLogValue() slog.Value {
	return slog.GroupValue(
		slog.String("host", a.Host),
		slog.String("port", a.Port),
	)
}

// AsLogValue represents Server struct as slog.Value
// Used for logging
func (s *Server) AsLogValue() slog.Value {
	return slog.GroupValue(
		slog.String("host", s.Host),
		slog.String("port", s.Port),
		slog.Duration("timeout", s.Timeout),
		slog.Duration("idleTimeout", s.IdleTimeout),
	)
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

	use := os.Getenv("USE_TEST_SERVER")
	if use == "true" {
		cfg.UseTestApi = true
	} else {
		cfg.UseTestApi = false
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
