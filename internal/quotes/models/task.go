package models

import (
	"time"
)

const (
	TaskStatusSuccess    = "success"
	TaskStatusFail       = "fail"
	TaskStatusInProgress = "in_progress"
)

type TaskDTO struct {
	CurrencyPair
	TaskID int64      `db:"id"`
	Time   *time.Time `db:"time"`
	Ratio  *float64   `db:"ratio"`
	Status string     `db:"status"`
}

func (t TaskDTO) ToQuoteData() QuoteData {
	quote := QuoteData{}
	if t.Time != nil {
		quote.Time = *t.Time
	}
	if t.Ratio != nil {
		quote.Ratio = *t.Ratio
	}
	return quote
}
