package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

type Config struct {
	DB        DBConfig
	JWT       JWTConfig
	Port      int
	RoomAdress string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func (c DBConfig) Validate() error {
	if c.Host == "" {
		return errors.New("DB_HOST is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return errors.New("DB_PORT must be between 1 and 65535")
	}
	if c.User == "" {
		return errors.New("DB_USER is required")
	}
	if c.Password == "" {
		return errors.New("DB_PASSWORD is required")
	}
	if c.Name == "" {
		return errors.New("DB_NAME is required")
	}
	return nil
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name)
}

type JWTConfig struct {
	Secret           string
	AccessExpiry     time.Duration
}

func (c JWTConfig) Validate() error {
	if c.Secret == "" {
		return errors.New("JWT_SECRET is required")
	}
	if len(c.Secret) < 32 {
		return errors.New("JWT_SECRET must be at least 32 characters")
	}

	return nil
}


func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	_ = godotenv.Overload(".env.secrets")

	cfg := &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "room"),
		},
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", ""),
			AccessExpiry:     getEnvDuration("JWT_ACCESS_EXPIRY", 7*24*time.Hour),
		},
		Port: getEnvInt("PORT", 8083),
		RoomAdress: getEnv("ROOM_SERVICE_ADDR", ":9090"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if err := c.DB.Validate(); err != nil {
		return fmt.Errorf("database: %w", err)
	}
	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("jwt: %w", err)
	}
	return nil
}
