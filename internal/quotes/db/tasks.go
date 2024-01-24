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

type taskStorage struct {
	conn   *sqlx.DB
	logger *slog.Logger
}

func NewTaskStorage(conn *sqlx.DB, logger *slog.Logger) *taskStorage {
	return &taskStorage{conn: conn, logger: logger}
}

func (t taskStorage) Save(ctx context.Context, pair models.CurrencyPair) (int64, error) {
	var id int64
	err := t.conn.QueryRowContext(ctx, "INSERT INTO refresh_task(base, counter) values($1, $2) RETURNING id;", pair.Base, pair.Counter).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("cannot create task: %w", errors.Join(models.ErrInsertErr, err))
	}
	return id, nil
}

func (t taskStorage) Get(ctx context.Context, id int64) (models.TaskDTO, error) {
	result := models.TaskDTO{}
	row := t.conn.QueryRowContext(ctx, "SELECT  base, counter, ratio, time, finished_at FROM refresh_task WHERE id = $1;", id)
	err := row.Scan(&result.Base, &result.Counter, &result.Ratio, &result.Time, &result.TimeFinished)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.TaskDTO{}, fmt.Errorf("cannot get quote for task: %w", models.ErrNoRows)
	case err != nil:
		return models.TaskDTO{}, err
	}

	return result, nil
}
