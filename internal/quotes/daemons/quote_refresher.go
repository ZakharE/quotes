package daemons

import (
	"context"
	"errors"
	"log/slog"
	"plata_card_quotes/internal/quotes/config"
	"plata_card_quotes/internal/quotes/models"
	"sync"
	"time"
)

type RefreshTaskRepository interface {
	GetUnprocessed(ctx context.Context, limit int) ([]models.TaskDTO, error)
	MarkSuccessAndUpdate(ctx context.Context, task models.Quote, taskIds []int64) error
	MarkFailed(ctx context.Context, ids []int64) error
}

type CurrencyQuotesClient interface {
	GetQuote(ctx context.Context, pair models.CurrencyPair) (models.Quote, error)
}

type notificationEventDaemon struct {
	cfg                   *config.DaemonSettings
	logger                *slog.Logger
	refreshTaskRepository RefreshTaskRepository
	quotesClient          CurrencyQuotesClient
}

func NewQuoteRefresherDaemon(
	cfg *config.DaemonSettings,
	quotesClient CurrencyQuotesClient,
	refreshTaskRepository RefreshTaskRepository,
	logger *slog.Logger,
) *notificationEventDaemon {
	return &notificationEventDaemon{
		cfg:                   cfg,
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
	wg := sync.WaitGroup{}
	for pair := range pairsByIds {
		wg.Add(1)
		ids := pairsByIds[pair]
		go func(pair models.CurrencyPair, ids []int64) {
			defer wg.Done()
			quote, err := n.quotesClient.GetQuote(ctx, pair)
			if err != nil {
				n.logger.Error("Cannot get quote for pair", "pair", pair, "error", err)
				n.logger.Info("Mark tasks as 'failed'")
				err = n.refreshTaskRepository.MarkFailed(ctx, ids)
				if err != nil {
					n.logger.Error("Cannot mark tasks as failed", "ids", ids, "error", err)
				}
			}

			err = n.refreshTaskRepository.MarkSuccessAndUpdate(ctx, quote, ids)
			if err != nil {
				n.logger.Error("Cannot save quote in database", "pair", pair, "error", err)
			}

			n.logger.Info("Tasks was updated", "ids", ids)
		}(pair, ids)
	}
	wg.Wait()
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
	return n.cfg.BatchSize
}

func (n notificationEventDaemon) BatchSleep() time.Duration {
	return time.Second * n.cfg.BatchSleep
}

func (n notificationEventDaemon) NoWorkSleep() time.Duration {
	return time.Second * n.cfg.NoWorkSleep
}
func (n notificationEventDaemon) Name() string {
	return "task_refresher"
}
