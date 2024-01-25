package server

import (
	"context"
	"errors"
	"log/slog"
	"plata_card_quotes/internal/quotes/models"
	"plata_card_quotes/internal/quotes/service"

	"net/http"
	"os"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
)

type quotesServer struct {
	logger        *slog.Logger
	mux           *chi.Mux
	quotesService *service.QuotesService
}

func NewQuotesServer(
	logger *slog.Logger,
	mux *chi.Mux,
	quotesService *service.QuotesService,
) *quotesServer {
	return &quotesServer{
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

	s := &http.Server{
		Handler: qs.mux,
		Addr:    ":8080",
	}

	err = s.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		qs.logger.Warn("Some error occured during shutdown", "error", err)
		os.Exit(1)
	}
}

func (qs *quotesServer) RefreshQuote(ctx context.Context, request RefreshQuoteRequestObject) (RefreshQuoteResponseObject, error) {
	body := request.Body
	if body.Base == body.Counter {
		return RefreshQuotedefaultJSONResponse{
			Body:       models.Error{Message: "Fields 'from' and 'to must differ!"},
			StatusCode: 400,
		}, nil
	}
	id, err := qs.quotesService.CreateRefreshTask(ctx, body.ToCurrencyPair())
	if err != nil {
		return RefreshQuotedefaultJSONResponse{models.Error{
			Message: "Something went wrong",
		}, 500}, nil
	}
	return RefreshQuote200JSONResponse(models.RefreshTask{Id: &id}), nil
}

func (qs *quotesServer) GetLastQuote(ctx context.Context, request GetLastQuoteRequestObject) (GetLastQuoteResponseObject, error) {
	pair := models.CurrencyPair{
		Base:    request.Params.BaseCurrency,
		Counter: request.Params.CounterCurrency,
	}
	quote, err := qs.quotesService.GetLastQuote(ctx, pair)
	if errors.Is(err, models.ErrNoRows) {
		return GetLastQuote404JSONResponse(models.Error{Message: "no task with such id"}), nil
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
	if err != nil {
		return GetTaskdefaultJSONResponse{models.Error{Message: "something went wrong"}, 500}, nil
	}
	return GetTask200JSONResponse(task.ToQuote()), nil
}
