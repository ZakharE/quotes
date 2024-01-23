package main

import (
	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"plata_card_quotes/internal/quotes/clients/currency_api"
	"plata_card_quotes/internal/quotes/config"
	"plata_card_quotes/internal/quotes/db"
	"plata_card_quotes/internal/quotes/server"
	"plata_card_quotes/internal/quotes/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs/")
	err := viper.ReadInConfig() // Find and read the config file

	if err != nil {
		logger.Error("Cannot find config file", "error", err)
		os.Exit(1)
	}

	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		logger.Error("Cannot unmarshal config file", "error", err)
		os.Exit(1)
	}

	client := currency_api.NewCurrencyQuotesClient(logger)
	conn := sqlx.MustConnect("postgres", cfg.Db.ConnectionString())
	storage := db.NewTaskStorage(conn, logger)

	srv := server.NewQuotesServer(logger,
		chi.NewRouter(),
		service.NewQuotesService(client, storage),
	)
	srv.Start()
}
