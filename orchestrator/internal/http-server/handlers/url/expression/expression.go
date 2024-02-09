package expression

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"orchestrator/internal/lib/api/response"
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
// @Tags Expressions
// @Accept json
// @Produce json
// @Param body body ExpressionRequest true "Запрос на создание выражения"
// @Success 201 {object} ExpressionResponse
// @Router /expression [post]
func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", "handlers.url.expression.New"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty")

			render.JSON(w, r, response.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", err.Error())
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", err.Error())
			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		expression := req.Expression

		// TODO add service layer call

		log.Info("expression calculating", slog.String("expression", expression))

		responseOK(w, r, "1")
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id string) {
	render.JSON(w, r, Response{
		Id:       id,
		Response: response.OK(),
	})
}
