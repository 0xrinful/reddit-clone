package background

import (
	"fmt"
	"log/slog"
	"sync"
)

type Worker struct {
	wg sync.WaitGroup
}

func New() *Worker {
	return &Worker{}
}

func (w *Worker) Run(fn func()) {
	w.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("background task panic", "error", fmt.Sprintf("%v", err))
			}
		}()

		fn()
	})
}

func (w *Worker) Wait() {
	w.wg.Wait()
}
