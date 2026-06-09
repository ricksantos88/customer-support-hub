package config

import (
	"fmt"
	"os"
	"strings"
	"time"

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

	AuthAccessTokenTTLMinutes int
	AuthRefreshTokenTTLHours  int
	AuthSessionTTLHours       int
	AuthCacheTTLMinutes       int
	AuthRateLimitPerMinute    int

	CORSAllowedOrigins string
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
	v.SetDefault("DB_PASSWORD", "support123") //TODO: change me in the future
	v.SetDefault("DB_NAME", "customer_support")
	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("AUTH_ACCESS_TOKEN_TTL_MINUTES", 15)
	v.SetDefault("AUTH_REFRESH_TOKEN_TTL_HOURS", 24)
	v.SetDefault("AUTH_SESSION_TTL_HOURS", 24)
	v.SetDefault("AUTH_CACHE_TTL_MINUTES", 30)
	v.SetDefault("AUTH_RATE_LIMIT_PER_MINUTE", 60)
	v.SetDefault("CORS_ALLOWED_ORIGINS", "*")

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
		AppName:                   v.GetString("APP_NAME"),
		Environment:               v.GetString("APP_ENV"),
		Host:                      v.GetString("APP_HOST"),
		Port:                      v.GetString("APP_PORT"),
		LogLevel:                  strings.ToUpper(v.GetString("APP_LOG_LEVEL")),
		DatabaseURL:               dbURL,
		RedisAddr:                 fmt.Sprintf("%s:%s", v.GetString("REDIS_HOST"), v.GetString("REDIS_PORT")),
		RedisDB:                   v.GetInt("REDIS_DB"),
		RedisPass:                 v.GetString("REDIS_PASSWORD"),
		JWTSecret:                 v.GetString("JWT_SECRET"),
		AuthAccessTokenTTLMinutes: v.GetInt("AUTH_ACCESS_TOKEN_TTL_MINUTES"),
		AuthRefreshTokenTTLHours:  v.GetInt("AUTH_REFRESH_TOKEN_TTL_HOURS"),
		AuthSessionTTLHours:       v.GetInt("AUTH_SESSION_TTL_HOURS"),
		AuthCacheTTLMinutes:       v.GetInt("AUTH_CACHE_TTL_MINUTES"),
		AuthRateLimitPerMinute:    v.GetInt("AUTH_RATE_LIMIT_PER_MINUTE"),
		CORSAllowedOrigins:        v.GetString("CORS_ALLOWED_ORIGINS"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	if cfg.AuthAccessTokenTTLMinutes <= 0 {
		cfg.AuthAccessTokenTTLMinutes = int((15 * time.Minute).Minutes())
	}
	if cfg.AuthRefreshTokenTTLHours <= 0 {
		cfg.AuthRefreshTokenTTLHours = 24
	}
	if cfg.AuthSessionTTLHours <= 0 {
		cfg.AuthSessionTTLHours = 24
	}
	if cfg.AuthCacheTTLMinutes <= 0 {
		cfg.AuthCacheTTLMinutes = 30
	}
	if cfg.AuthRateLimitPerMinute <= 0 {
		cfg.AuthRateLimitPerMinute = 60
	}

	return cfg, nil
}
