package db

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"plata_card_quotes/internal/quotes/models"
)

type taskStorage struct {
	conn   *sqlx.DB
	logger *slog.Logger
}

func NewTaskStorage(conn *sqlx.DB, logger *slog.Logger) *taskStorage {
	return &taskStorage{conn: conn, logger: logger}
}

func (t taskStorage) Save(ctx context.Context, pair models.CurrencyPair) (int, error) {
	//TODO implement me
	panic("implement me")
}
