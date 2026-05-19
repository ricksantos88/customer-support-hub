package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppName     string
	Environment string
	Host        string
	Port        string
	LogLevel    string

	DatabaseURL string
	RedisAddr   string
	RedisDB     int
	RedisPass   string

	JWTSecret string
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("APP_NAME", "customer-support-hub")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_HOST", "0.0.0.0")
	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("APP_LOG_LEVEL", "INFO")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_USER", "support")
	v.SetDefault("DB_PASSWORD", "support123")
	v.SetDefault("DB_NAME", "customer_support")
	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("REDIS_DB", 0)

	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = "development"
	}

	v.SetConfigFile(fmt.Sprintf(".env.%s", environment))
	_ = v.ReadInConfig()
	v.SetConfigFile(".env")
	_ = v.MergeInConfig()

	dbURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		v.GetString("DB_HOST"),
		v.GetString("DB_PORT"),
		v.GetString("DB_USER"),
		v.GetString("DB_PASSWORD"),
		v.GetString("DB_NAME"),
		v.GetString("DB_SSL_MODE"),
	)

	cfg := &Config{
		AppName:     v.GetString("APP_NAME"),
		Environment: v.GetString("APP_ENV"),
		Host:        v.GetString("APP_HOST"),
		Port:        v.GetString("APP_PORT"),
		LogLevel:    strings.ToUpper(v.GetString("APP_LOG_LEVEL")),
		DatabaseURL: dbURL,
		RedisAddr:   fmt.Sprintf("%s:%s", v.GetString("REDIS_HOST"), v.GetString("REDIS_PORT")),
		RedisDB:     v.GetInt("REDIS_DB"),
		RedisPass:   v.GetString("REDIS_PASSWORD"),
		JWTSecret:   v.GetString("JWT_SECRET"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}
