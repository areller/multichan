package multichan

import (
	"sync"
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

func syncMapLength(m *sync.Map) int {
	c := 0
	m.Range(func (k, v interface{}) bool {
		c++
		return true
	})
	return c
}