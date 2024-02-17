package execution_time

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net/http"
	"orchestrator/internal/app"
	"orchestrator/internal/lib/api/response"
	"orchestrator/internal/services/settings_service"
)

type Request struct {
	ExecutionTime int    `json:"execution_time" validate:"required" example:"1"`
	OperationType string `json:"operation_type" validate:"required,eq=plus|eq=minus|eq=mult|eq=div" example:"plus"`
}

type Response struct {
	Id       string            `json:"id"`
	Response response.Response `json:"response"`
}

// @Summary Установка нового времени выполнения
// @Description operation_type: minus, plus, mult, div. execution_time > 0
// @Tags Settings
// @Accept json
// @Produce json
// @Param body body Request true "Установка времени выполнения"
// @Success 201 {object} Response
// @Router /settings/execution-time [post]
func New(log *slog.Logger, application *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		logger := log.With(
			slog.String("op", "handlers.settings_service.execution_time.New"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
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

		err = application.SettingService.UpdateSettingsExecutionTime(ctx, req.OperationType, req.ExecutionTime)
		if err != nil {
			if errors.Is(err, settings_service.ErrValidationOperationTime) {
				logger.Info("validation of execution operation time failed", err.Error())
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error("Invalid execution time"))
				return
			}
			logger.Error("internal error", err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Internal Error"))
			return
		}
		logger.Info("expression calculating", slog.Int("expression", req.ExecutionTime))

		responseOK(w, r, uuid.New().String())
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id string) {
	render.JSON(w, r, Response{
		Id:       id,
		Response: response.OK(),
	})
}
