package models

import (
	"strings"
	"time"
)

type TaskDTO struct {
	CurrencyPair
	TaskID       int64      `db:"id"`
	Time         *time.Time `db:"time"`
	Ratio        *float64   `db:"ratio"`
	TimeFinished *time.Time `db:"time_finished"`
}

func (t TaskDTO) ToQuote() Quote {
	quote := Quote{}
	if t.Time != nil {
		quote.Time = t.Time.String()
	}
	if t.Ratio != nil {
		quote.Ratio = *t.Ratio
	}
	return quote
}

type CurrencyPair struct {
	Base    Currency `db:"base"`
	Counter Currency `db:"counter"`
}

func (rt NewRefreshTask) ToCurrencyPair() CurrencyPair {
	return CurrencyPair{
		Base:    Currency(strings.ToLower(string(rt.Base))),
		Counter: Currency(strings.ToLower(string(rt.Counter))),
	}
}
