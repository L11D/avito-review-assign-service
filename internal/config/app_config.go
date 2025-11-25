package config

import (
	"errors"
	"log/slog"
	"os"
	"time"
)

const (
	DEFAULT_SHUTDOWN_TIMEOUT    = 10 * time.Second
	DEFAULT_READ_HEADER_TIMEOUT = 10 * time.Second
)

type Config struct {
	DBUser            string
	DBPassword        string
	DBHost            string
	DBPort            string
	DBName            string
	HTTPPort          string
	ShutdownTimeout   time.Duration
	ReadHeaderTimeout time.Duration
}

func getEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return val, errors.New("ENV " + key + " is missing")
	}

	return val, nil
}

func getEnvOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		slog.Warn("ENV " + key + " is missing, using default " + def)

		return def
	}

	return val
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		slog.Warn("ENV " + key + " is missing, using default " + def.String())

		return def
	}

	duration, err := time.ParseDuration(val)
	if err != nil {
		slog.Warn("ENV " + key + " is invalid, using default " + def.String())

		return def
	}

	return duration
}

func LoadConfig() (*Config, error) {
	dbUser, err := getEnv("POSTGRES_USER")
	if err != nil {
		return nil, err
	}

	dbPass, err := getEnv("POSTGRES_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := getEnv("POSTGRES_DB")
	if err != nil {
		return nil, err
	}

	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")

	httpPort := getEnvOrDefault("HTTP_PORT", "8080")
	shutdownTimeout := getEnvDuration("SHUTDOWN_TIMEOUT", DEFAULT_SHUTDOWN_TIMEOUT)

	return &Config{
		DBUser:            dbUser,
		DBPassword:        dbPass,
		DBHost:            dbHost,
		DBPort:            dbPort,
		DBName:            dbName,
		HTTPPort:          httpPort,
		ShutdownTimeout:   shutdownTimeout,
		ReadHeaderTimeout: DEFAULT_READ_HEADER_TIMEOUT,
	}, nil
}

func (c *Config) GetDBSource() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" +
		c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}
