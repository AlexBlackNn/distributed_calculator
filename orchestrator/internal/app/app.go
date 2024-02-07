package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/config"
	"sso/internal/services/auth_service"
	patroni "sso/storage/patroni"
	redis "sso/storage/redis-sentinel"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	//init storage
	// TODO: seems to need factory here
	//storage, err := postgres.New(cfg.StoragePath) //uncomment for postgres
	storage, err := patroni.New(cfg)

	if err != nil {
		panic(err)
	}
	//init cache
	tokenCache := redis.New(cfg)

	//init auth_service service (auth_service)
	authService := auth_service.New(log, storage, tokenCache, cfg)
	grpcApp := grpcapp.New(log, authService, cfg)
	return &App{
		GRPCSrv: grpcApp,
	}
}
