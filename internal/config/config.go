package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config holds application configuration settings
type Config struct {
	Server ServerConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host string
	Port string
}

// LoadConfig initializes and returns application configuration
func LoadConfig() (*Config, error) {
	// Initialize with defaults
	cfg := &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: "8080",
		},
	}

	// Try to load from .env file first
	if err := loadEnvFile(".env"); err != nil {
		// Only return error if it's not a "file not found" error
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
		// If file doesn't exist, that's ok - continue with defaults and env vars
	}

	// Override with environment variables
	applyEnvironment(cfg)

	return cfg, nil
}

// loadEnvFile reads environment variables from a file
func loadEnvFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("malformed line in %s: %s", filename, line)
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading %s: %w", filename, err)
	}

	return nil
}

// applyEnvironment overrides config with environment variables
func applyEnvironment(cfg *Config) {
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	}
}
