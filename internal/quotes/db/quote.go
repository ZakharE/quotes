package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"plata_card_quotes/internal/quotes/models"
)

type quoteStorage struct {
	conn   *sqlx.DB
	logger *slog.Logger
}

func NewQuoteStorage(conn *sqlx.DB, logger *slog.Logger) *quoteStorage {
	return &quoteStorage{conn: conn, logger: logger}
}
func (qs quoteStorage) Get(ctx context.Context, pair models.CurrencyPair) (models.QuoteData, error) {
	var quote models.QuoteData
	err := qs.conn.GetContext(ctx, &quote, "SELECT ratio, time  FROM quote WHERE base = $1 AND counter = $2;", pair.Base, pair.Counter)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.QuoteData{}, fmt.Errorf("cannot get quote : %w", models.ErrNoRows)
	case err != nil:
		return models.QuoteData{}, err
	}

	return quote, nil
}
