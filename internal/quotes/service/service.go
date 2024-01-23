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
	client CurrencyQuotesClient
}

func NewQuotesService(client CurrencyQuotesClient) *QuotesService {
	return &QuotesService{client: client}
}

func (qs QuotesService) CreateRefreshTask(ctx context.Context, pair models.CurrencyPair) (models.Quote, error) {

	return qs.client.GetQuote(ctx, pair)
}
