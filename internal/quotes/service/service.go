package service

import (
	"context"
	"log/slog"
	"plata_card_quotes/internal/quotes/models"
)

type CurrencyQuotesClient interface {
	GetQuote(ctx context.Context, pair models.CurrencyPair) (models.Quote, error)
}

type QuoteRepository interface {
	Get(ctx context.Context, pair models.CurrencyPair) (models.QuoteData, error)
}

type RefreshTaskRepository interface {
	Save(ctx context.Context, pair models.CurrencyPair) (int64, error)
	Get(ctx context.Context, id int64) (models.TaskDTO, error)
}

type QuotesService struct {
	logger                *slog.Logger
	client                CurrencyQuotesClient
	refreshTaskRepository RefreshTaskRepository
	quoteRepository       QuoteRepository
}

func NewQuotesService(
	logger *slog.Logger,
	client CurrencyQuotesClient,
	refreshTaskRepository RefreshTaskRepository,
	quoteRepository QuoteRepository,
) *QuotesService {
	return &QuotesService{
		logger:                logger,
		client:                client,
		refreshTaskRepository: refreshTaskRepository,
		quoteRepository:       quoteRepository,
	}
}

func (qs QuotesService) CreateRefreshTask(ctx context.Context, pair models.CurrencyPair) (int64, error) {
	id, err := qs.refreshTaskRepository.Save(ctx, pair)
	if err != nil {
		qs.logger.Error("cannot create new task", "error", err)
	}
	return id, err
}

func (qs QuotesService) GetTask(ctx context.Context, id int64) (models.TaskDTO, error) {
	res, err := qs.refreshTaskRepository.Get(ctx, id)
	if err != nil {
		qs.logger.Error("cannot retrieve task", "error", err)
	}
	return res, err
}

func (qs QuotesService) GetLastQuote(ctx context.Context, pair models.CurrencyPair) (models.QuoteData, error) {
	quote, err := qs.quoteRepository.Get(ctx, pair)
	if err != nil {
		qs.logger.Error("cannot retrieve last quote", "error", err)
	}
	return quote, err
}
