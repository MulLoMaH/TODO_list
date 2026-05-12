package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MulLoMaH/TODO_list.git/internal/core/domain"
	core_logger "github.com/MulLoMaH/TODO_list.git/internal/core/logger"
	core_http_request "github.com/MulLoMaH/TODO_list.git/internal/core/transport/http/request"
	core_http_response "github.com/MulLoMaH/TODO_list.git/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"`
	TasksAvarageCompletionTime *string  `json:"tasks_avarage_completion_time"`
}

// GetStatistics    godoc
// @Summary         Получение статистики
// @Discriotion     Получение статистики по задачам с опциональной фильтрацией по user_id и/или временному промежутку
// @Tags            statistics
// @Produce         json
// @Param           user_id query int false "Фильтрация статистики по конкретному пользователю"
// @Param           from query string false "Начало промежутка рассмотрения статистики (включительно) формат: YYYY-MM-DD"
// @Param           to query string false  "Конец промежутка рассмотрения статистики (не включительно) формат: YYYY-MM-DD"
// @Success         200 {object} GetStatisticsResponse "Успешное получение статистики"
// @Failure 		400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router          /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/from/to query params",
		)

		return
	}

	sratistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)

		return
	}

	response := toDTOFromDomain(sratistics)

	responseHandler.JSONResponse(response, http.StatusOK)
}

// type queryParams struct {
// 	userID *int
// 	from   *time.Time
// 	to     time.Time
// }

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAvarageCompletionTime != nil {
		duration := statistics.TasksAvarageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAvarageCompletionTime: avgTime,
	}
}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := core_http_request.GetIntQueryParams(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParams(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParams(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userID, from, to, nil
}
