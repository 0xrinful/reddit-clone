package background

import (
	"fmt"
	"log/slog"
	"sync"
)

type Worker struct {
	wg     sync.WaitGroup
	logger *slog.Logger
}

func New(logger *slog.Logger) *Worker {
	return &Worker{logger: logger}
}

func (w *Worker) Run(fn func()) {
	w.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				w.logger.Error("background task panic", "error", fmt.Sprintf("%v", err))
			}
		}()

		fn()
	})
}

func (w *Worker) Wait() {
	w.wg.Wait()
}
