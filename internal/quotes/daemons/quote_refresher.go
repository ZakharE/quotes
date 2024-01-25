package daemons

import (
	"context"
	"errors"
	"log/slog"
	"plata_card_quotes/internal/quotes/models"
	"time"
)

type RefreshTaskRepository interface {
	GetUnprocessed(ctx context.Context, limit int) ([]models.TaskDTO, error)
	SetProcessedAndUpdated(ctx context.Context, task models.TaskDTO, taskIds []int64) error
}

type CurrencyQuotesClient interface {
	GetQuote(ctx context.Context, pair models.CurrencyPair) (models.TaskDTO, error)
}

type notificationEventDaemon struct {
	logger                *slog.Logger
	refreshTaskRepository RefreshTaskRepository
	quotesClient          CurrencyQuotesClient
}

func NewQuoteRefresherDaemon(
	quotesClient CurrencyQuotesClient,
	refreshTaskRepository RefreshTaskRepository,
	logger *slog.Logger,
) *notificationEventDaemon {
	return &notificationEventDaemon{
		quotesClient:          quotesClient,
		refreshTaskRepository: refreshTaskRepository,
		logger:                logger,
	}
}

func (n notificationEventDaemon) ProcessBatch(ctx context.Context, batchSize int) error {
	tasks, err := n.refreshTaskRepository.GetUnprocessed(ctx, batchSize)
	if errors.Is(err, models.ErrNoRows) {
		return ErrNoWork
	}

	pairsByIds := splitByCurrency(tasks)
	for pair := range pairsByIds {
		quote, err := n.quotesClient.GetQuote(ctx, pair)
		if err != nil {
			n.logger.Error("Cannot get tasks for pair", "pair", pair, "error", err)
			continue
		}
		// update by ids and with new values
		err = n.refreshTaskRepository.SetProcessedAndUpdated(ctx, quote, pairsByIds[pair])
		if err != nil {
			n.logger.Error("Cannot save quote in database", "pair", pair, "error", err)
			continue
		}

	}

	return nil
}

func splitByCurrency(tasks []models.TaskDTO) map[models.CurrencyPair][]int64 {
	result := make(map[models.CurrencyPair][]int64, len(tasks)/2) // maybe preallocate capacity
	for i := range tasks {
		ids, ok := result[tasks[i].CurrencyPair]
		if !ok {
			sl := make([]int64, 0, len(tasks))
			sl = append(sl, tasks[i].TaskID)
			result[tasks[i].CurrencyPair] = sl
			continue
		}
		result[tasks[i].CurrencyPair] = append(ids, tasks[i].TaskID)
	}
	return result
}

func (n notificationEventDaemon) BatchSize() int {
	return 100
}

func (n notificationEventDaemon) BatchSleep() time.Duration {
	return time.Second * 10
}

func (n notificationEventDaemon) NoWorkSleep() time.Duration {
	return time.Second * 20
}