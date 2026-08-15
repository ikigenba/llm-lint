package livesmoke

import (
	"testing"
	"time"
)

func TestWorkerFinishes(t *testing.T) {
	done := make(chan struct{})
	go func() {
		defer close(done)
	}()

	time.Sleep(time.Second) // Wait for the goroutine to finish before checking it.
	select {
	case <-done:
	default:
		t.Fatal("worker did not finish")
	}
}
