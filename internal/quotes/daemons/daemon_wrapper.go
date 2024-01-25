package daemons

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

var (
	ErrNoWork = errors.New("no new work for daemon")
)

type Daemon interface {
	ProcessBatch(ctx context.Context, batchSize int) error
	BatchSize() int
	BatchSleep() time.Duration
	NoWorkSleep() time.Duration
}

type multiDaemonWrapper struct {
	daemons []Daemon
	logger  *slog.Logger
}

func NewMultiDaemonWrapper(logger *slog.Logger) *multiDaemonWrapper {
	return &multiDaemonWrapper{logger: logger}
}

func (mdw *multiDaemonWrapper) Register(d Daemon) {
	mdw.daemons = append(mdw.daemons, d)
}

func (mdw *multiDaemonWrapper) Start(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(len(mdw.daemons))
	for _, d := range mdw.daemons {
		go func(d Daemon) {
			defer wg.Done()
			mdw.startDaemon(ctx, d)
		}(d)
	}
	wg.Wait()
}

func (mdw *multiDaemonWrapper) startDaemon(ctx context.Context, d Daemon) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := d.ProcessBatch(ctx, d.BatchSize())
			switch {
			case errors.Is(err, ErrNoWork):
				{
					mdw.logger.Info("no work. sleep")
					time.Sleep(d.NoWorkSleep())
				}
			case err != nil:
				mdw.logger.Error("err during batch processing", "error", err)
			}
			time.Sleep(d.BatchSleep())
		}
	}
}
