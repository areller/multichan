package multichan

import (
	"testing"
)

func BenchmarkSimple(b *testing.B) {
	c := New()
	lis := c.Listen()

	for n := 0; n < b.N; n++ {
		c.Input() <- n
		<- lis.Output()
	}
}

func BenchmarkPubSub(b *testing.B) {
	c := New()
	lisA := c.Listen()
	lisB := c.Listen()

	drain := func (l *Listener) {
		for {
			<- l.Output()
		}
	}
	go drain(lisA)
	go drain(lisB)

	for n := 0; n < b.N; n++ {
		c.Input() <- n
	}
}

func BenchmarkBufferedChannel(b *testing.B) {
	c := NewWithBuffer(b.N)
	for n := 0; n < b.N; n++ {
		c.Input() <- n
	}
}

func BenchmarkBufferedChannelWithListener(b *testing.B) {
	c := NewWithBuffer(b.N)
	lis := c.Listen()
	for n := 0; n < b.N; n++ {
		c.Input() <- n
		<- lis.Output()
	}
}

func BenchmarkInfiniteChannel(b *testing.B) {
	c := NewInfinite()
	for n := 0; n < b.N; n++ {
		c.Input() <- n
	}
}

func BenchmarkInfiniteChannelWithListener(b *testing.B) {
	c := NewInfinite()
	lis := c.Listen()
	for n := 0; n < b.N; n++ {
		c.Input() <- n
		<- lis.Output()
	}
}