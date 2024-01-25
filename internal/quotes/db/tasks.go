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
	err := t.conn.QueryRowContext(ctx, "INSERT INTO refresh_task(base, counter) VALUES($1, $2) RETURNING id;", pair.Base, pair.Counter).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("cannot create task: %w", errors.Join(models.ErrInsertErr, err))
	}
	return id, nil
}

func (t taskStorage) Get(ctx context.Context, id int64) (models.TaskDTO, error) {
	result := models.TaskDTO{}
	row := t.conn.QueryRowContext(ctx, "SELECT  base, counter, ratio, time, is_finished FROM refresh_task WHERE id = $1;", id)
	err := row.Scan(&result.Base, &result.Counter, &result.Ratio, &result.Time, &result.TimeFinished, &result.IsFinished)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.TaskDTO{}, fmt.Errorf("cannot get quote for task: %w", models.ErrNoRows)
	case err != nil:
		return models.TaskDTO{}, err
	}

	return result, nil
}

func (t taskStorage) GetUnprocessed(ctx context.Context, limit int) ([]models.TaskDTO, error) {
	result := make([]models.TaskDTO, 0, limit)
	err := t.conn.SelectContext(ctx, &result, "SELECT  id, base, counter FROM refresh_task WHERE is_finished = FALSE;")
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: no new tasks were found ", models.ErrNoRows)
	}
	return result, nil
}

func (t taskStorage) SetProcessedAndUpdated(ctx context.Context, quote models.Quote, taskIds []int64) error {
	tx, err := t.conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot start transaction: %w: %w", models.ErrTransaction, err)
	}

	query, args, err := sqlx.In("UPDATE refresh_task SET ratio = ?, time = ?, is_finished = TRUE WHERE id IN (?);", quote.Ratio, quote.Time, taskIds)
	if err != nil {
		rlErr := tx.Rollback()
		if rlErr != nil {
			return fmt.Errorf("cannot rollback transaction after query rebinding: %w", errors.Join(models.ErrTransaction, err, rlErr))
		}
		return fmt.Errorf("unable to update refresh_task table, %w", errors.Join(models.ErrTransactionRollback, err))
	}

	query = tx.Rebind(query)
	_, err = tx.ExecContext(ctx, query, args...)

	if err != nil {
		rlErr := tx.Rollback()
		if rlErr != nil {
			return fmt.Errorf("cannot rollback transaction: %w", errors.Join(models.ErrTransaction, err, rlErr))
		}
		return fmt.Errorf("unable to rebind query, %w", errors.Join(models.ErrTransactionRollback, err))
	}
	query = `INSERT INTO quote (base, counter, ratio, time)
											  VALUES (:base, :counter, :ratio, :time)
											  ON CONFLICT (base, counter) DO UPDATE SET ratio=:ratio, time=:time;`

	_, err = tx.NamedExecContext(ctx, query, quote)

	if err != nil {
		rlErr := tx.Rollback()
		if rlErr != nil {
			return fmt.Errorf("cannot rollback transaction: %w", errors.Join(models.ErrTransaction, err, rlErr))
		}
		return fmt.Errorf("unable to update quote table, %w", errors.Join(models.ErrTransactionRollback, err))
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cannot commit transaction: %w: %w", models.ErrTransaction, err)
	}
	return err
}
