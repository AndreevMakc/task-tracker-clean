package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTP   HTTPConfig
	PG     PGConfig
	App    AppConfig
	TG     TGConfig
}

type HTTPConfig struct {
	Port        string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

type PGConfig struct {
	URL          string
	MaxPoolSize  int
	ConnAttempts int
	ConnTimeout  time.Duration
}

type TGConfig struct {
	BotToken string
}

type AppConfig struct {
	Name string
}

func New() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Port:        getEnv("HTTP_PORT", "8080"),
			Timeout:     getDurationEnv("HTTP_TIMEOUT", 5*time.Second),
			IdleTimeout: getDurationEnv("HTTP_IDLE_TIMEOUT", 60*time.Second),
		},
		PG: PGConfig{
			URL:          getEnv("PG_URL", "postgres://postgres:postgres@localhost:5432/tasktracker?sslmode=disable"),
			MaxPoolSize:  getIntEnv("PG_MAX_POOL_SIZE", 10),
			ConnAttempts: getIntEnv("PG_CONN_ATTEMPTS", 10),
			ConnTimeout:  getDurationEnv("PG_CONN_TIMEOUT", time.Second),
		},
		App: AppConfig{
			Name: getEnv("APP_NAME", "task-tracker-clean"),
		},
		TG: TGConfig{
			BotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
		log.Printf("invalid value for %s: %v, using default", key, err)
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return duration
		}
		log.Printf("invalid value for %s: %v, using default", key, err)
	}
	return defaultValue
}
