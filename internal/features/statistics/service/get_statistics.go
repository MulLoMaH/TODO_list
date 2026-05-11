package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/MulLoMaH/TODO_list.git/internal/core/domain"
	core_errors "github.com/MulLoMaH/TODO_list.git/internal/core/errors"
)

//Проверка на соответствие интерфейсу
// var _ statistics_transport_http.StatisticsService = (*StatisticsService)(nil)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"'to' must be after 'from': %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	tasks, err := s.statisticsRepository.GetStatistics(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf(
			"get tasks from retository: %w",
			err,
		)
	}

	statistics := calcStatistics(tasks)

	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.NewStatistics(0, 0, nil, nil)
	}

	tasksCreated := len(tasks)

	tasksCompleted := 0
	var totalCompletionDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		completionDuration := task.ComplelionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAvarageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)

		tasksAvarageCompletionTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		&tasksCompletedRate,
		tasksAvarageCompletionTime,
	)
}
