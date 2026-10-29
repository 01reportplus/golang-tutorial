package pipeline

import (
	"context"
	"errors"
	"slices"
	"sync/atomic"
	"testing"
	"time"
)

func TestMapOrderAndBound(t *testing.T) {
	started := make(chan struct{}, 2)
	release := make(chan struct{})
	var active, peak atomic.Int64
	done := make(chan struct{})
	var got []int
	var err error
	go func() {
		defer close(done)
		got, err = Map(context.Background(), 2, []int{1, 2, 3, 4}, func(_ context.Context, n int) (int, error) {
			current := active.Add(1)
			defer active.Add(-1)
			for old := peak.Load(); current > old; old = peak.Load() {
				if peak.CompareAndSwap(old, current) {
					break
				}
			}
			if n <= 2 {
				started <- struct{}{}
				<-release
			}
			return n * n, nil
		})
	}()
	for range 2 {
		select {
		case <-started:
		case <-time.After(3 * time.Second):
			close(release)
			t.Fatal("workers did not start")
		}
	}
	close(release)
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("workers did not finish")
	}
	if err != nil || !slices.Equal(got, []int{1, 4, 9, 16}) || peak.Load() != 2 || active.Load() != 0 {
		t.Fatalf("got=%v err=%v peak=%d active=%d", got, err, peak.Load(), active.Load())
	}
}

func TestMapCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	started := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		_, err := Map(ctx, 1, []int{1, 2}, func(ctx context.Context, _ int) (int, error) {
			close(started)
			<-ctx.Done()
			return 0, ctx.Err()
		})
		done <- err
	}()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("worker did not start")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("got %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("cancellation leaked worker")
	}
}

func TestMapErrorAndInvalidInput(t *testing.T) {
	want := errors.New("work failed")
	fn := func(context.Context, int) (int, error) { return 0, want }
	if _, err := Map(context.Background(), 2, []int{1, 2}, fn); !errors.Is(err, want) {
		t.Fatalf("got %v, want %v", err, want)
	}
	if _, err := Map(context.Background(), 0, []int{1}, fn); err == nil {
		t.Fatal("expected invalid worker error")
	}
	if _, err := Map[int, int](context.Background(), 1, nil, nil); err == nil {
		t.Fatal("expected nil function error")
	}
	got, err := Map(context.Background(), 2, []int{}, fn)
	if err != nil || len(got) != 0 {
		t.Fatalf("empty result=%v err=%v", got, err)
	}
}
