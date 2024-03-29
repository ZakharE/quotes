package quotes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"path"
	"plata_card_quotes/internal/quotes/clients/currency"
	"plata_card_quotes/internal/quotes/config"
	"plata_card_quotes/internal/quotes/daemons"
	"plata_card_quotes/internal/quotes/db"
	"plata_card_quotes/internal/quotes/server"
	"plata_card_quotes/internal/quotes/service"
	"strings"
)

func StartApp() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	viper.SetEnvPrefix("ENV")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.MustBindEnv("db.host", "DB_HOST")
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs/")
	err := viper.ReadInConfig() // Find and read the config file
	viper.Debug()

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
	connectionString := cfg.Db.ConnectionString()
	conn := sqlx.MustConnect("postgres", connectionString)
	currentWd, err := os.Getwd()

	if err != nil {
		return
	}

	fs := os.DirFS(path.Join(currentWd, "sql"))
	goose.SetBaseFS(fs)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(conn.DB, "quotes"); err != nil {
		panic(err)
	}

	taskStorage := db.NewTaskStorage(conn, logger)
	quoteStorage := db.NewQuoteStorage(conn, logger)

	client := currency.NewCurrencyQuotesClient(logger)

	daemonWrapper := daemons.NewMultiDaemonWrapper(logger)
	daemon := daemons.NewQuoteRefresherDaemon(&cfg.DaemonSettings.TaskRefresher, client, taskStorage, logger)
	daemonWrapper.Register(daemon)

	srv := server.NewQuotesServer(
		&cfg.Server,
		logger,
		chi.NewRouter(),
		service.NewQuotesService(logger, client, taskStorage, quoteStorage),
		daemonWrapper,
	)

	srv.Start()
}
