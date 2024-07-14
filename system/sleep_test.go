package system

import (
	"context"
	"testing"
	"time"
)

func TestSleepWithContext(t *testing.T) {
	t.Run("Sleeps for the specified duration", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		startTime := time.Now()
		SleepWithContext(ctx, 50*time.Millisecond)
		elapsedTime := time.Since(startTime)

		if elapsedTime < 50*time.Millisecond {
			t.Errorf("SleepWithContext didn't sleep for the expected duration")
		}
	})

	t.Run("Returns immediately if context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel the context immediately

		startTime := time.Now()
		SleepWithContext(ctx, 50*time.Millisecond)
		elapsedTime := time.Since(startTime)

		if elapsedTime > 10*time.Millisecond {
			t.Errorf("SleepWithContext should return immediately when the context is canceled")
		}
	})
}
