package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/signal"
	"plata_card_quotes/internal/quotes/config"
	"plata_card_quotes/internal/quotes/models"
	"plata_card_quotes/internal/quotes/service"
	"syscall"
	"time"

	"net/http"
	"os"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
)

type quotesServer struct {
	cfg           *config.Server
	srv           *http.Server
	logger        *slog.Logger
	mux           *chi.Mux
	quotesService *service.QuotesService
}

func NewQuotesServer(
	cfg *config.Server,
	logger *slog.Logger,
	mux *chi.Mux,
	quotesService *service.QuotesService,
) *quotesServer {
	return &quotesServer{
		cfg:           cfg,
		logger:        logger,
		mux:           mux,
		quotesService: quotesService,
	}
}

func (qs *quotesServer) Start() {
	swagger, err := GetSwagger()
	if err != nil {
		qs.logger.Error("unable to get swagger.exit")
		os.Exit(1)
	}

	swagger.Servers = nil
	qs.mux.Use(middleware.OapiRequestValidator(swagger))

	h := NewStrictHandler(qs, nil)
	HandlerFromMux(h, qs.mux)

	qs.srv = &http.Server{
		Handler: qs.mux,
		Addr:    fmt.Sprintf(":%d", qs.cfg.Port),
	}

	serverCtx := qs.gracefulShutdown()
	err = qs.srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		qs.logger.Warn("Some error occured during shutdown", "error", err)
		os.Exit(1)
	}
	qs.logger.Debug("Shutting down the server. Wait...")
	<-serverCtx.Done()
	qs.logger.Debug("Server was shut down")
}

func (qs *quotesServer) RefreshQuote(ctx context.Context, request RefreshQuoteRequestObject) (RefreshQuoteResponseObject, error) {
	pair := models.NewCurrencyPair(request.Body.Base, request.Body.Counter)
	if !models.IsCurrencySupported(pair.Base) {
		return RefreshQuote400JSONResponse(models.Error{Message: fmt.Sprintf("Currency '%s' is not supported", pair.Base)}), nil
	}
	if !models.IsCurrencySupported(pair.Counter) {
		return RefreshQuote400JSONResponse(models.Error{Message: fmt.Sprintf("Currency '%s' is not supported", pair.Counter)}), nil
	}
	if pair.Base == pair.Counter {
		return RefreshQuote400JSONResponse(models.Error{Message: "Fields 'from' and 'to must differ!"}), nil
	}
	id, err := qs.quotesService.CreateRefreshTask(ctx, pair)
	if err != nil {
		return RefreshQuotedefaultJSONResponse{models.Error{
			Message: "Something went wrong",
		}, 500}, nil
	}
	return RefreshQuote200JSONResponse(models.RefreshTask{Id: &id}), nil
}

func (qs *quotesServer) GetLastQuote(ctx context.Context, request GetLastQuoteRequestObject) (GetLastQuoteResponseObject, error) {
	pair := models.NewCurrencyPair(request.Params.BaseCurrency, request.Params.CounterCurrency)
	if !models.IsCurrencySupported(pair.Base) {
		return GetLastQuote400JSONResponse(models.Error{Message: fmt.Sprintf("Currency '%s' is not supported", pair.Base)}), nil
	}
	if !models.IsCurrencySupported(pair.Counter) {
		return GetLastQuote400JSONResponse(models.Error{Message: fmt.Sprintf("Currency '%s' is not supported", pair.Counter)}), nil
	}
	quote, err := qs.quotesService.GetLastQuote(ctx, pair)
	if errors.Is(err, models.ErrNoRows) {
		return GetLastQuote404JSONResponse(models.Error{Message: "no quote with such currency pair"}), nil
	}
	if err != nil {
		return GetLastQuotedefaultJSONResponse{models.Error{Message: "something went wrong"}, 500}, nil
	}
	return GetLastQuote200JSONResponse(quote), nil
}

func (qs *quotesServer) GetTask(ctx context.Context, request GetTaskRequestObject) (GetTaskResponseObject, error) {
	task, err := qs.quotesService.GetTask(ctx, request.Id)
	if errors.Is(err, models.ErrNoRows) {
		return GetTask404JSONResponse(models.Error{Message: "no task with such id"}), nil
	}

	if task.Status != models.TaskStatusSuccess {
		return GetTask425JSONResponse(models.TaskResponseError{Message: "task still in progress. try again later", Status: task.Status}), nil
	}

	if err != nil {
		return GetTaskdefaultJSONResponse{models.Error{Message: "something went wrong"}, 500}, nil
	}
	return GetTask200JSONResponse(task.ToQuoteData()), nil
}

func (qs *quotesServer) gracefulShutdown() context.Context {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	go func() {
		<-signalCh
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				qs.logger.Error("graceful shutdown timed out.. forcing exit.")
				os.Exit(1)
			}
		}()

		// Trigger graceful shutdown
		err := qs.srv.Shutdown(shutdownCtx)
		if err != nil {
			qs.logger.Error("graceful shutdown timed out.. forcing exit.")
			os.Exit(1)
		}
		serverStopCtx()
	}()
	return serverCtx
}
