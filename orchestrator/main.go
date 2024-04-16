package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log/slog"
	"net/http"
	_ "orchestrator/docs"
	"orchestrator/internal/app"
	"orchestrator/internal/config"
	"orchestrator/internal/http-server/handlers/calculation/expression"
	"orchestrator/internal/http-server/handlers/calculation/operations"
	"orchestrator/internal/http-server/handlers/calculation/result"
	"orchestrator/internal/http-server/handlers/monitoring/worker"
	"orchestrator/internal/http-server/handlers/settings/execution_time"
	projectLogger "orchestrator/internal/http-server/middleware/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title           Swagger API
// @version         1.0
// @description     This is a distributed calculation server.
// @contact.name   API Support
// @license.name  Apache 2.0
// @license.calculation   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
//
//go:generate go run github.com/swaggo/swag/cmd/swag init
func main() {
	// init config
	cfg := config.MustLoad()
	// init logger
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))

	// init app
	ctx := context.Background()
	application := app.New(log, cfg)
	go application.OrchestrationService.ParseResponse(ctx)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(projectLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/expression", func(r chi.Router) {
		r.Get("/{uuid}", result.New(log, application))
		r.Post("/", expression.New(log, application))
	})

	router.Route("/", func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		))
	})

	router.Route("/monitoring", func(r chi.Router) {
		r.Get("/worker", worker.New(log, application, cfg))
	})

	router.Route("/settings", func(r chi.Router) {
		r.Post("/execution-time", execution_time.New(log, application))
	})

	router.Route("/operations", func(r chi.Router) {
		r.Get("/", operations.GetOperationsWithPaginationHandler(log, application))
	})

	// graceful stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      router,
		ReadTimeout:  time.Duration(10 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(10 * time.Second),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	signalType := <-stop
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
