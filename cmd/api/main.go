package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ricksantos88/customer-support-hub/internal/application/auth"
	"github.com/ricksantos88/customer-support-hub/internal/config"
	"github.com/ricksantos88/customer-support-hub/internal/database"
	"github.com/ricksantos88/customer-support-hub/internal/infrastructure/cache"
	httpiface "github.com/ricksantos88/customer-support-hub/internal/interfaces/http"
	handlers "github.com/ricksantos88/customer-support-hub/internal/interfaces/http/handlers"
	"github.com/ricksantos88/customer-support-hub/internal/interfaces/http/middleware"
	"github.com/ricksantos88/customer-support-hub/internal/logger"
	postgresrepo "github.com/ricksantos88/customer-support-hub/internal/repositories/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.LogLevel)
	slog.SetDefault(log)

	dbConn, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	redisClient := cache.NewRedisClient(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)
	if err := cache.Ping(context.Background(), redisClient); err != nil {
		slog.Warn("redis unavailable, continuing with postgres-backed sessions", "error", err)
	}

	agentRepo := postgresrepo.NewAgentRepository(dbConn)
	sessionRepo := postgresrepo.NewSessionRepository(dbConn)
	sessionCache := cache.NewSessionCache(redisClient)
	authService := auth.NewService(
		agentRepo,
		sessionRepo,
		sessionCache,
		cfg.JWTSecret,
		cfg.AppName,
		time.Duration(cfg.AuthAccessTokenTTLMinutes)*time.Minute,
		time.Duration(cfg.AuthRefreshTokenTTLHours)*time.Hour,
		time.Duration(cfg.AuthSessionTTLHours)*time.Hour,
	)
	authHandler := handlers.NewAuthHandler(authService)
	authMiddleware := middleware.NewBearerAuthMiddleware(cfg.JWTSecret, sessionRepo)

	app := httpiface.NewRouter(httpiface.RouterDependencies{
		AuthHandler:        authHandler,
		AuthMiddleware:     authMiddleware,
		IPRateLimiter:      middleware.NewIPRateLimiter(cfg.AuthRateLimitPerMinute),
		AgentRateLimiter:   middleware.NewAgentRateLimiter(cfg.AuthRateLimitPerMinute),
		CORSAllowedOrigins: cfg.CORSAllowedOrigins,
	})
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	slog.Info("starting api", "env", cfg.Environment, "addr", addr)
	if err := app.Listen(addr); err != nil {
		slog.Error("api stopped", "error", err)
	}
}
