package multichan

import (
	"time"
)

func tryWithTimeout(ts time.Duration, what func()) bool {
	done := make(chan struct{})
	go func() {
		what()
		close(done)
	}()

	select {
	case <-time.After(ts):
		return false
	case <-done:
		return true
	}
}