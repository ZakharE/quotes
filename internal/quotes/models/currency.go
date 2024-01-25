package models

var supportedCurrency = map[string]struct{}{
	"eur": {},
	"usd": {},
	"mxn": {},
}

func IsCurrencySupported(currency string) bool {
	_, ok := supportedCurrency[currency]
	return ok
}

type CurrencyPair struct {
	Base    string `db:"base"`
	Counter string `db:"counter"`
}
