package currency_api

import (
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"net/http"
	"plata_card_quotes/internal/quotes/models"
	"time"
)

var (
	base_api                  = "https://cdn.jsdelivr.net/gh/fawazahmed0"
	ErrNotFound               = errors.New("data not found")
	ErrNotSupportedStatusCode = errors.New("unexpected status from server")
)

type currenciesClient struct {
	client *resty.Client
	logger *slog.Logger
}

func NewCurrencyQuotesClient(logger *slog.Logger) *currenciesClient {
	client := resty.New()
	client.SetBaseURL(base_api)
	return &currenciesClient{
		client: client,
		logger: logger,
	}
}

func (c currenciesClient) GetQuote(ctx context.Context, pair models.CurrencyPair) (models.Quote, error) {
	url := "/currency-api@1/latest/currencies/{from}/{to}.json"
	type resp map[string]interface{}
	response, err := c.client.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
		},
		).
		SetResult(&resp{}).
		SetPathParams(map[string]string{
			"from": string(pair.From),
			"to":   string(pair.To),
		}).
		Get(url)

	if err != nil {
		c.logger.Error("cannot get quote", "error", err)
		return models.Quote{}, err
	}

	statusCode := response.StatusCode()
	switch statusCode {
	case http.StatusOK:
		result := response.Result().(*resp)
		quote, err := c.parseResponse(pair, *result)
		if err != nil {
			c.logger.Error(
				"cannot parse response",
				"error", err,
				"resp", result,
			)
			return models.Quote{}, errors.New("cannot parse response")
		}
		return quote, nil
	case http.StatusNotFound:
		return models.Quote{}, ErrNotFound
	default:
		c.logger.Warn(
			"cannot parse response",
			"error", err,
			"status_code", statusCode,
			"body", response.RawResponse,
		)
		return models.Quote{}, ErrNotSupportedStatusCode
	}
}

func (c currenciesClient) parseResponse(pair models.CurrencyPair, response map[string]interface{}) (models.Quote, error) {
	date, err := time.Parse(time.DateOnly, response["date"].(string))

	if err != nil {
		return models.Quote{}, err
	}
	val := response[string(pair.To)]
	if val == nil {
		return models.Quote{}, err
	}

	return models.Quote{
		Pair:  pair,
		Date:  date,
		Value: val.(float64),
	}, err
}
