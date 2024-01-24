package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"path"
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

	getwd, err := os.Getwd()
	if err != nil {
		return
	}
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
	// --- setup database
	conn := sqlx.MustConnect("postgres", cfg.Db.ConnectionString())

	fs := os.DirFS(path.Join(getwd, "sql"))
	goose.SetBaseFS(fs)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(conn.DB, "quotes"); err != nil {
		panic(err)
	}

	taskStorage := db.NewTaskStorage(conn, logger)
	quoteStorage := db.NewQuoteStorage(conn, logger)

	client := currency_api.NewCurrencyQuotesClient(logger)
	srv := server.NewQuotesServer(logger,
		chi.NewRouter(),
		service.NewQuotesService(logger, client, taskStorage, quoteStorage),
	)

	srv.Start()
}
