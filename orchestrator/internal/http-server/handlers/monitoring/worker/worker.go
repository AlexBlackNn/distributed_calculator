package worker

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"orchestrator/internal/app"
	"orchestrator/internal/lib/api/response"
)

type Response struct {
	Id       string            `json:"id"`
	Response response.Response `json:"response"`
}

// @Summary Получение количества активных воркеров
// @Description Получает количество воркеров доступных для выполнения задачи
// @Tags Monitoring
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /monitoring/worker [get]
func New(log *slog.Logger, application *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		log := log.With(
			slog.String("op", "handlers.monitoring.worker.New"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		result, err := application.MonitoringService.GetActiveWorkers(ctx)
		// TODO: think about this error
		if err != nil {
			log.Error("some errors", err.Error())
		}
		log.Info("expression result", slog.Float64("result", result))
		responseOK(w, r, uuid.New().String(), result)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id string, result float64) {
	render.JSON(w, r, Response{
		Id:       id,
		Response: response.Result(result),
	})
}
