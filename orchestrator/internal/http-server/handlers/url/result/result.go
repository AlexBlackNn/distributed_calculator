package result

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"orchestrator/internal/app"
	"orchestrator/internal/lib/api/response"
)

type Response struct {
	Id       string            `json:"id"`
	Response response.Response `json:"response"`
}

// @Summary Получение результата по идентификатору
// @Description Получает результат по указанному идентификатору из хранилища
// @Tags Results
// @Accept json
// @Produce json
// @Param uid path string true "Идентификатор результата"
// @Success 200 {object} Response
// @Router /{uid} [get]
func New(log *slog.Logger, application *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		log := log.With(
			slog.String("op", "handlers.url.expression.New"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		uid := chi.URLParam(r, "uid")
		if uid == "" {
			log.Info("absent uid")
			render.JSON(w, r, response.Error("absent uid"))
		}

		result, err := application.OrchestrationService.CalculationResult(ctx, uid)
		// TODO: think about this error
		if err != nil {
			log.Error("some errors", err.Error())
		}
		log.Info("expression result", slog.Float64("result", result))

		responseOK(w, r, uid, result)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id string, result float64) {
	render.JSON(w, r, Response{
		Id:       id,
		Response: response.Result(result),
	})
}
