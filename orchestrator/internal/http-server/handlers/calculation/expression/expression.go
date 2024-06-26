package expression

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
	"orchestrator/internal/domain/models"
	"orchestrator/internal/http-server/handlers/utils"
	"orchestrator/internal/lib/api/response"
	"strings"
)

type Request struct {
	Expression string `json:"expression" validate:"required"`
}

type Response struct {
	Id       string            `json:"id"`
	Response response.Response `json:"response"`
}

// @Summary Создание нового выражения
// @Description Создает новое выражение на сервере
// @Tags Calculations
// @Accept json
// @Produce json
// @Param body body Request true "Запрос на создание выражения"
// @Success 201 {object} Response
// @Router /expression [post]
// @Security BearerAuth
func New(log *slog.Logger, application *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		tokenString := r.Header.Get("Authorization")
		token := strings.TrimPrefix(tokenString, "Bearer")
		if !utils.JWTCheck(token) {
			log.Error("jwt token check failed")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("bad jwt token"))
			return
		}
		uid, name, err := utils.JWTParse(token)
		if err != nil {
			log.Error("jwt parsing failed")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("jwt parsing failed"))
			return
		}
		appUser := models.User{uid, name}
		ctx := context.Background()
		log := log.With(
			slog.String("op", "handlers.calculation.expression.New"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			render.Status(r, http.StatusBadRequest)
			log.Error("invalid request", err.Error())
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		id, err := application.OrchestrationService.CalculationRequest(ctx, req.Expression, appUser)
		if err != nil {
			log.Error("invalid request", err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Internal Error"))
			return
		}
		log.Info("expression calculating", slog.String("expression", req.Expression))

		responseOK(w, r, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id string) {
	render.JSON(w, r, Response{
		Id:       id,
		Response: response.OK(),
	})
}
