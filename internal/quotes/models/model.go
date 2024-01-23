package models

import (
	"strings"
	"time"
)

type Quote struct {
	Pair  CurrencyPair
	Date  time.Time
	Value float64 //would be better if we can store in struct with separate decimal and floating part.
}

type CurrencyPair struct {
	From Currency
	To   Currency
}

func (rt NewRefreshTask) ToCurrencyPair() CurrencyPair {
	return CurrencyPair{
		From: Currency(strings.ToLower(string(rt.From))),
		To:   Currency(strings.ToLower(string(rt.To))),
	}
}
