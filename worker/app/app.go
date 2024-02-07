package app

import (
	"context"
	"log/slog"
	"worker/internal/config"
	service "worker/internal/service"
	"worker/internal/transport/rabbit"
)

type MessageBrokerInterface interface {
	Send(context.Context, any, *config.Config)
	Receive(*config.Config)
	Stop()
}

type App struct {
	MessageBroker MessageBrokerInterface
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	// init service layer
	calculator := service.New()
	//init rabbitmq
	messageBroker, err := rabbit.New(log, cfg, calculator)
	if err != nil {
		if err != nil {
			panic(err)
		}
	}

	return &App{
		MessageBroker: messageBroker,
	}
}
