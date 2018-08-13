package multichan

import (
	"sync"
	"time"
)

type simpleTuple struct {
	a interface{}
	b bool
}

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

func sendWithTimeout(ts time.Duration, c chan interface{}, msg interface{}) bool {
	select {
	case c <- msg:
		return true
	case <-time.After(ts):
		return false
	}
}

func recvWithTimeout(ts time.Duration, c chan interface{}) (interface{}, bool) {
	select {
	case msg := <-c:
		return msg, true
	case <-time.After(ts):
		return nil, false
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

func makeST(a interface{}, b bool) simpleTuple {
	return simpleTuple{
		a: a,
		b: b,
	}
}