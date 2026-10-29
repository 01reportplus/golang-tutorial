// Package pipeline demonstrates bounded parallel work with ordered output.
package pipeline

import (
	"context"
	"errors"
	"sync"
)

// Map applies fn with at most workers active calls and preserves input order.
// The first error cancels remaining work. No partial results are returned.
// fn must respect ctx and must not panic; Map waits for every started worker.
func Map[T, R any](parent context.Context, workers int, inputs []T, fn func(context.Context, T) (R, error)) ([]R, error) {
	if workers < 1 {
		return nil, errors.New("workers must be positive")
	}
	if fn == nil {
		return nil, errors.New("function must not be nil")
	}
	ctx, cancel := context.WithCancelCause(parent)
	defer cancel(nil)
	out := make([]R, len(inputs))
	jobs := make(chan int)
	var wg sync.WaitGroup
	for range min(workers, len(inputs)) {
		wg.Go(func() {
			for index := range jobs {
				if ctx.Err() != nil {
					return
				}
				result, err := fn(ctx, inputs[index])
				if err != nil {
					cancel(err)
					return
				}
				out[index] = result // Every index has exactly one writer.
			}
		})
	}
feed:
	for index := range inputs {
		select {
		case jobs <- index:
		case <-ctx.Done():
			break feed
		}
	}
	close(jobs) // The only sender owns closure.
	wg.Wait()
	if err := context.Cause(ctx); err != nil {
		return nil, err
	}
	return out, nil
}
