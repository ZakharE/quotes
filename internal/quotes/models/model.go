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
	IsFinished   bool       `db:"is_finished"`
}

type Quote struct {
	CurrencyPair
	QuoteData
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

type CurrencyPair struct {
	Base    string `db:"base"`
	Counter string `db:"counter"`
}

func (rt NewRefreshTask) ToCurrencyPair() CurrencyPair {
	return CurrencyPair{
		Base:    strings.ToLower(rt.Base),
		Counter: strings.ToLower(rt.Counter),
	}
}
