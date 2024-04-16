package user_operations

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"orchestrator/internal/app"
	"orchestrator/internal/domain/models"
	"orchestrator/internal/http-server/handlers/utils"
	"orchestrator/internal/lib/api/response"
	"strconv"
)

type Response struct {
	Id       string            `json:"id"`
	Response response.Response `json:"response"`
}

// @Summary Получение операций с пагинацией пользователя
// @Description Переход с 1 страницы на случайную не предусмотрен! Пагинация быстрая с поиском по индексу. В качестве курсора пустое значение для начала, потом скопировать ПОСЛЕДНЮЮ дату ПОЛЯ CreatedAt , например 2024-02-18T16:27:05.271813Z
// @Tags Operations
// @Accept json
// @Produce json
// @Param page_size query int false "Размер страницы" default(2)
// @Param cursor query string false "Курсор для пагинации"
// @Success 200 {array} Response
// @Router /operations/user [get]
// @Security BearerAuth
func GetUserOperationsWithPaginationHandler(log *slog.Logger, application *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, name, err := utils.ParseJWTToken(r)
		appUser := models.User{uid, name}
		if err != nil {
			log.Error("jwt token check failed")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("bad jwt token"))
			return
		}

		ctx := context.Background()
		log := log.With(
			slog.String("op", "handlers.operations.GetOperationsWithPaginationHandler"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		pageSizeStr := r.URL.Query().Get("page_size")
		cursor := r.URL.Query().Get("cursor")
		// TODO: move to settings
		pageSize := 10 // Default page size
		// TODO: move all validators to validator folder in utils
		if pageSizeStr != "" {
			pageSize, err = strconv.Atoi(pageSizeStr)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error("Invalid page size"))
				return
			}
		}
		operations, err := application.OrchestrationService.GetUserOperationsWithPagination(ctx, pageSize, cursor, appUser)
		if err != nil {
			log.Error("Failed to get operations with pagination: ", err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Internal server error"))
			return
		}
		render.JSON(w, r, operations)
	}
}
