package health

import (
	"context"
	"testing"
	"time"
)

func TestCheckerCollectsConcurrentResultsWithoutDeadlock(t *testing.T) {
	checker := Checker{Timeout: time.Second, Checks: map[string]Check{
		"database": func(context.Context) error { return nil },
		"worker":   func(context.Context) error { return nil },
	}}
	done := make(chan []Result, 1)
	go func() { results, _ := checker.Run(context.Background()); done <- results }()
	select {
	case results := <-done:
		if len(results) != 2 {
			t.Fatalf("result count=%d", len(results))
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("health checks deadlocked while publishing results")
	}
}
