package monitoring_service

import (
	"context"
	"errors"
	"github.com/levigross/grequests"
	"log/slog"
	"orchestrator/internal/config"
	"orchestrator/message_broker"
)

type MonitoringService struct {
	log           *slog.Logger
	messageBroker message_broker.MessageBrokerInterface
	cfg           *config.Config
}

// New returns a new instance of MonitoringService
func New(
	log *slog.Logger,
	cfg *config.Config,
) *MonitoringService {
	return &MonitoringService{
		log: log,
		cfg: cfg,
	}
}

func (ms *MonitoringService) GetActiveWorkers(
	ctx context.Context,
) (float64, error) {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: monitoring_service.GetActiveWorkers"),
	)
	log.Info("check if operation was calculated")
	// TODO: get from config
	rabbitUrl := "http://guest:guest@localhost:15672/api/queues/%2f/operation"
	// Send a GET request to the RabbitMQ Management API to get queue details
	//TODO: change for client with timeout
	resp, err := grequests.Get(rabbitUrl, nil)
	if err != nil {
		return 0, err
	}
	var data map[string]interface{}
	err = resp.JSON(&data)
	if err != nil {
		return 0, err
	}

	if consumers, ok := data["consumers"].(float64); ok {
		return consumers, nil
	} else {
		// TODO: create error for service
		return 0, errors.New("Consumers information not available for this queue.")
	}

}
