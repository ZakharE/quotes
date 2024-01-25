package models

const (
	CurrencyEUR = "eur"
	CurrencyUSD = "usd"
	CurrencyMXN = "mxn"
)

var supportedCurrency = map[string]struct{}{
	CurrencyEUR: {},
	CurrencyUSD: {},
	CurrencyMXN: {},
}

func IsCurrencySupported(currency string) bool {
	_, ok := supportedCurrency[currency]
	return ok
}

type CurrencyPair struct {
	Base    string `db:"base"`
	Counter string `db:"counter"`
}
