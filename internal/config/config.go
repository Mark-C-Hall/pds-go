package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}

	Database struct {
		ConnectionString string
		MaxOpenConns     int
		MaxIdleConns     int
		ConnMaxLifetime  time.Duration
	}

	Log struct {
		Level string
	}
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{}

	// Server Setings
	port, err := requireEnv("SERVER_PORT")
	if err != nil {
		return nil, err
	}
	readTimeout, err := requireDuration("SERVER_READ_TIMEOUT")
	if err != nil {
		return nil, err
	}
	writeTimeout, err := requireDuration("SERVER_WRITE_TIMEOUT")
	if err != nil {
		return nil, err
	}
	idleTimeout, err := requireDuration("SERVER_IDLE_TIMEOUT")
	if err != nil {
		return nil, err
	}
	cfg.Server.Port = port
	cfg.Server.ReadTimeout = readTimeout
	cfg.Server.WriteTimeout = writeTimeout
	cfg.Server.IdleTimeout = idleTimeout

	// Database Settings
	connectionString, err := requireEnv("DB_CONNECTION_STRING")
	if err != nil {
		return nil, err
	}
	maxOpenConns, err := requireInt("DB_MAX_OPEN_CONNS")
	if err != nil {
		return nil, err
	}
	maxIdleConns, err := requireInt("DB_MAX_IDLE_CONNS")
	if err != nil {
		return nil, err
	}
	connMaxLifetime, err := requireDuration("DB_CONN_MAX_LIFETIME")
	if err != nil {
		return nil, err
	}
	cfg.Database.ConnectionString = connectionString
	cfg.Database.MaxOpenConns = maxOpenConns
	cfg.Database.MaxIdleConns = maxIdleConns
	cfg.Database.ConnMaxLifetime = connMaxLifetime

	// Log Settings
	logLevel, err := requireEnv("LOG_LEVEL")
	if err != nil {
		return nil, err
	}
	cfg.Log.Level = logLevel

	return cfg, nil
}

func requireEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("%s environment variable not defined", key)
	}
	return value, nil
}

func requireDuration(key string) (time.Duration, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return 0, fmt.Errorf("%s environment variable not defined", key)
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid duration value for %s: %w", key, err)
	}
	return duration, nil
}

func requireInt(key string) (int, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return 0, fmt.Errorf("%s environment variable not defined", key)
	}
	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer value for %s: %w", key, err)
	}
	return num, nil
}
