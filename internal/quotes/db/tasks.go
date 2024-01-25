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
	row := t.conn.QueryRowContext(ctx, "SELECT  base, counter, ratio, time, status FROM refresh_task WHERE id = $1;", id)
	err := row.Scan(&result.Base, &result.Counter, &result.Ratio, &result.Time, &result.Status)
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
	query := `SELECT  id, base, counter FROM refresh_task WHERE status = $1 AND last_attempt_at < NOW() - INTERVAL '1 minutes' LIMIT  $2;`
	err := t.conn.SelectContext(ctx, &result, query, models.TaskStatusInProgress, limit)

	if err != nil {
		return nil, fmt.Errorf("cannot retrieve rows: %w", err)
	}

	if len(result) == 0 {
		return nil, models.ErrNoRows
	}

	return result, nil
}

func (t taskStorage) MarkSuccessAndUpdate(ctx context.Context, quote models.Quote, taskIds []int64) error {
	return t.withTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		query, args, err := sqlx.In("UPDATE refresh_task SET ratio = ?, time = ?, status = ?, last_attempt_at = NOW() WHERE id IN (?);", quote.Ratio, quote.Time, models.TaskStatusSuccess, taskIds)
		if err != nil {
			return fmt.Errorf("unable to transform an update query, %w", errors.Join(models.ErrTransactionRollback, err))
		}

		query = tx.Rebind(query)
		_, err = tx.ExecContext(ctx, query, args...)

		if err != nil {
			return fmt.Errorf("unable to update resfresh task, %w", errors.Join(models.ErrTransactionRollback, err))
		}

		query = `
				 INSERT INTO quote (base, counter, ratio, time)
				 VALUES (:base, :counter, :ratio, :time)
				 ON CONFLICT (base, counter) DO UPDATE SET ratio=:ratio,
																		 time=:time;
				 `

		_, err = tx.NamedExecContext(ctx, query, quote)

		if err != nil {
			return fmt.Errorf("unable to update quote table, %w", errors.Join(models.ErrTransactionRollback, err))
		}
		return nil
	})

}

func (t taskStorage) MarkFailed(ctx context.Context, taskIds []int64) error {
	query, args, err := sqlx.In("UPDATE refresh_task SET  status = ?, last_attempt_at = NOW() WHERE id IN (?);", models.TaskStatusFail, taskIds)
	if err != nil {
		return fmt.Errorf("unable to transform an update query, %w", errors.Join(models.ErrTransactionRollback, err))
	}

	query = t.conn.Rebind(query)
	_, err = t.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("cannot mark tasks as failed: %w", err)
	}
	return nil
}

func (t taskStorage) withTx(ctx context.Context, f func(ctx context.Context, tx *sqlx.Tx) error) error {
	tx, err := t.conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot start transaction: %w", errors.Join(models.ErrTransaction, err))
	}
	err = f(ctx, tx)
	if err != nil {
		rlErr := tx.Rollback()
		if rlErr != nil {
			return fmt.Errorf("cannot rollback transaction: %w", errors.Join(models.ErrTransaction, err, rlErr))
		}
		return fmt.Errorf("cannot execute transaction: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cannot commit transaction: %w: %w", models.ErrTransaction, err)
	}
	return nil

}
