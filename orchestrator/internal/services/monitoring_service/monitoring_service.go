package monitoring_service

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"orchestrator/internal/config"
	"orchestrator/message_broker"
	"time"
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
	cfg *config.Config,
) (float64, error) {
	log := ms.log.With(
		slog.String("info", "SERVICE LAYER: monitoring_service.GetActiveWorkers"),
	)
	log.Info("starts getting active workers")

	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(cfg.RabbitUrlWorker)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}

	if consumers, ok := data["consumers"].(float64); ok {
		return consumers, nil
	} else {
		return 0, ErrNoWorkerInformation
	}

}
