package result

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"orchestrator/internal/app"
	"orchestrator/internal/lib/api/response"
	"orchestrator/internal/services/orchestrator_service"
)

type Response struct {
	Id       string            `json:"id"`
	Response response.Response `json:"response"`
}

// @Summary Получение результата по идентификатору
// @Description Получает результат по указанному идентификатору из хранилища
// @Tags Calculations
// @Accept json
// @Produce json
// @Param uid path string true "Идентификатор результата"
// @Success 200 {object} Response
// @Router /expression/{uuid} [get]
func New(log *slog.Logger, application *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		log := log.With(
			slog.String("op", "handlers.calculation.expression.New"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		currentUuid := chi.URLParam(r, "uuid")

		if !IsValidUUID(currentUuid) {
			log.Info("currentUuid", currentUuid)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("wrong uuid"))
			return
		}
		result, err := application.OrchestrationService.CalculationResult(ctx, currentUuid)

		if err != nil {
			log.Error("internal error", err.Error())
			if errors.Is(err, orchestrator_service.ErrNoOperation) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, response.Error("Operation with requested uuid not found"))
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Internal error"))
		}
		log.Info("expression result", slog.Float64("result", result))
		responseOK(w, r, currentUuid, result)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id string, result float64) {
	render.JSON(w, r, Response{
		Id:       id,
		Response: response.Result(result),
	})
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
