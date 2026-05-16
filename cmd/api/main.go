package main

import (
	"fmt"
	"log/slog"

	"github.com/ricksantos88/customer-support-hub/internal/config"
	httpiface "github.com/ricksantos88/customer-support-hub/internal/interfaces/http"
	"github.com/ricksantos88/customer-support-hub/internal/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.LogLevel)
	slog.SetDefault(log)

	app := httpiface.NewRouter()
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	slog.Info("starting api", "env", cfg.Environment, "addr", addr)
	if err := app.Listen(addr); err != nil {
		slog.Error("api stopped", "error", err)
	}
}
