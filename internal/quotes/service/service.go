package service

import (
	"context"
	"plata_card_quotes/internal/quotes/models"
)

type CurrencyQuotesClient interface {
	GetQuote(ctx context.Context, pair models.CurrencyPair) (models.Quote, error)
}

type RefreshTaskRepository interface {
	Save(ctx context.Context, pair models.CurrencyPair) (int, error)
}

type QuotesService struct {
	client                CurrencyQuotesClient
	refreshTaskRepository RefreshTaskRepository
}

func NewQuotesService(client CurrencyQuotesClient, refreshTaskRepository RefreshTaskRepository) *QuotesService {
	return &QuotesService{
		client:                client,
		refreshTaskRepository: refreshTaskRepository,
	}
}

func (qs QuotesService) CreateRefreshTask(ctx context.Context, pair models.CurrencyPair) (models.Quote, error) {

	return qs.client.GetQuote(ctx, pair)
}
