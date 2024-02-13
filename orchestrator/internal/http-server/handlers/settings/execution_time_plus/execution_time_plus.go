package execution_time_plus

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"orchestrator/internal/app"
	"orchestrator/internal/lib/api/response"
)

type Request struct {
	ExecutionTime int `json:"execution_time" validate:"required"`
}

type Response struct {
	Id       string            `json:"id"`
	Response response.Response `json:"response"`
}

// @Summary Создание нового выражения
// @Description Создает новое выражение на сервере
// @Tags Settings
// @Accept json
// @Produce json
// @Param body body Request true "Запрос на создание выражения"
// @Success 201 {object} Response
// @Router /settings/plus-execution-time [post]
func New(log *slog.Logger, application *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		logger := log.With(
			slog.String("op", "handlers.settings.execution_time.New"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			logger.Error("request body is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty request"))
			return
		}
		if err != nil {
			logger.Error("failed to decode request body", err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		logger.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			render.Status(r, http.StatusBadRequest)
			logger.Error("invalid request", err.Error())
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		err = application.SettingService.PlusExecutionTime(ctx, req.ExecutionTime)
		if err != nil {
			logger.Error("invalid request", err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Internal Error"))
		}
		logger.Info("expression calculating", slog.Int("expression", req.ExecutionTime))

		//TODO: change id
		responseOK(w, r, "1")
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id string) {
	render.JSON(w, r, Response{
		Id:       id,
		Response: response.OK(),
	})
}
