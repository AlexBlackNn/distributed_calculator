package app

import (
	"context"
	"distributed_calculator/worker/internal/config"
	service "distributed_calculator/worker/internal/service"
	"distributed_calculator/worker/internal/transport/rabbit"
	"log/slog"
)

type MessageBrokerInterface interface {
	Send(ctx context.Context, message any, string2 string) error
	Receive(string) error
	Stop() error
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
