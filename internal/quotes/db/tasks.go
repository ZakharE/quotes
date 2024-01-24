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

func (t taskStorage) Save(ctx context.Context, pair models.CurrencyPair) (int64, error) {
	result, err := t.conn.ExecContext(ctx, "INSERT INTO refresh_task(from, to) values($1, $2)", pair.From, pair.To)
	if err != nil {
		//TODO
		panic("implement wrapping")
	}
	id, err := result.LastInsertId()
	if err != nil {

		//TODO
		panic("implement wrapping")
	}
	return id, nil
}
