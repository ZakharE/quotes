package service

import "plata_card_quotes/internal/quotes/models"

type CurrencyQuotesClient interface {
	GetQuote(pair models.CurrencyPair) (models.Quote, error)
}

type QuotesService struct {
	client CurrencyQuotesClient
}

func NewQuotesService(client CurrencyQuotesClient) *QuotesService {
	return &QuotesService{client: client}
}
