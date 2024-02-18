package operations

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"orchestrator/internal/app"
	"orchestrator/internal/lib/api/response"
	"strconv"
)

type Response struct {
	Id       string            `json:"id"`
	Response response.Response `json:"response"`
}

// @Summary Получение операций с пагинацией
// @Description Переход с 1 страницы на случайную не предусмотрен! Пагинация быстрая с поиском по индексу. В качестве курсора пустое значение для начала, потом скопировать ПОСЛЕДНЮЮ дату ПОЛЯ CreatedAt , например 2024-02-18T16:27:05.271813Z
// @Tags Operations
// @Accept json
// @Produce json
// @Param page_size query int false "Размер страницы" default(2)
// @Param cursor query string false "Курсор для пагинации"
// @Success 200 {array} Response
// @Router /operations [get]
func GetOperationsWithPaginationHandler(log *slog.Logger, application *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		log := log.With(
			slog.String("op", "handlers.operations.GetOperationsWithPaginationHandler"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		pageSizeStr := r.URL.Query().Get("page_size")
		cursor := r.URL.Query().Get("cursor")

		pageSize := 10 // Default page size
		var err error
		if pageSizeStr != "" {
			pageSize, err = strconv.Atoi(pageSizeStr)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error("Invalid page size"))
				return
			}
		}

		operations, err := application.OrchestrationService.GetOperationsWithPagination(ctx, pageSize, cursor)
		if err != nil {
			log.Error("Failed to get operations with pagination: ", err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Internal server error"))
			return
		}

		render.JSON(w, r, operations)
	}
}
