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
func (qs quoteStorage) Get(ctx context.Context, pair models.CurrencyPair) (models.Quote, error) {
	var quote models.Quote
	row := qs.conn.QueryRowContext(ctx, "SELECT ratio, time  FROM quote WHERE base = $1 AND counter = $2;", pair.Base, pair.Counter)
	err := row.Scan(&quote.Time, &quote.Ratio)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Quote{}, fmt.Errorf("cannot get quote : %w", models.ErrNoRows)
	case err != nil:
		return models.Quote{}, err
	}

	return quote, nil
}
