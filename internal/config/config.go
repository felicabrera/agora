// Package config loads the API server's configuration from the environment.
//
// Environment rather than a file, for one reason: the values that will be added
// here as the system grows (database DSN, log signing keys) are secrets, and a
// config file is something people commit by accident.
package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config is the API server's runtime configuration.
type Config struct {
	// Addr is the listen address, e.g. ":8080".
	Addr string
	// ReadHeaderTimeout bounds how long a client may take to send its headers.
	// Without it a slowloris client can hold a connection open indefinitely.
	ReadHeaderTimeout time.Duration
	// ShutdownTimeout bounds graceful shutdown.
	ShutdownTimeout time.Duration
	// LogLevel is one of debug, info, warn, error.
	LogLevel string
}

// Load reads the configuration from the environment and applies defaults.
func Load() (*Config, error) {
	c := &Config{
		Addr:              env("AGORA_ADDR", ":8080"),
		ReadHeaderTimeout: envDuration("AGORA_READ_HEADER_TIMEOUT", 10*time.Second),
		ShutdownTimeout:   envDuration("AGORA_SHUTDOWN_TIMEOUT", 15*time.Second),
		LogLevel:          env("AGORA_LOG_LEVEL", "info"),
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) validate() error {
	switch strings.ToLower(c.LogLevel) {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("config: AGORA_LOG_LEVEL must be one of debug, info, warn, error; got %q", c.LogLevel)
	}
	if c.Addr == "" {
		return fmt.Errorf("config: AGORA_ADDR must not be empty")
	}
	return nil
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}
