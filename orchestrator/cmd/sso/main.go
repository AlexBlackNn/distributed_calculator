package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/tracing"
	"syscall"
)

func main() {
	// init config
	cfg := config.MustLoad()
	// init logger
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))
	// init app
	application := app.New(log, cfg)

	tp, err := tracing.Init("sso service", cfg)
	if err != nil {
		log.Error("init tracing error", err.Error())
	}
	defer func() {
		if err = tp.Shutdown(context.Background()); err != nil {
			log.Error("Error shutting down tracer provider", err.Error())
		}
	}()
	// start grpc_transport server
	go application.GRPCSrv.MustRun()
	// graceful stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	signalType := <-stop
	application.GRPCSrv.Stop()
	log.Info(
		"application stopped",
		slog.String("signalType",
			signalType.String()),
	)
}

// our environments
const (
	envLocal = "local"
	envDemo  = "demo"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{
					Level:     slog.LevelDebug,
					AddSource: true,
				},
			),
		)
	case envDemo:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{
					Level:     slog.LevelDebug,
					AddSource: true,
				},
			),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{
					Level:     slog.LevelInfo,
					AddSource: true,
				},
			),
		)
	}
	return log
}
