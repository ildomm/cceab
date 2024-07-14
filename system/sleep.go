package system

import (
	"context"
	"time"
)

// SleepWithContext sleeps for the specified duration, or until the context is canceled.
func SleepWithContext(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(duration):
	}
}
