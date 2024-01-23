package main

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"os"
	"plata_card_quotes/internal/quotes/clients/currency_api"
	"plata_card_quotes/internal/quotes/server"
	"plata_card_quotes/internal/quotes/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := currency_api.NewCurrencyQuotesClient(logger)
	srv := server.NewQuotesServer(logger,
		chi.NewRouter(),
		service.NewQuotesService(client),
	)
	srv.Start()
}
