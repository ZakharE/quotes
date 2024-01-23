package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"plata_card_quotes/internal/quotes/clients/currency_api"
	"plata_card_quotes/internal/quotes/models"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := currency_api.NewCurrencyQuotesClient(logger)
	res, err := client.GetQuote(models.CurrencyPair{
		From: models.CurrencyEUR,
		To:   models.CurrencyUSD,
	})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(res)
}
