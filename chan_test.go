package multichan

import (
	"time"
	"testing"
	"github.com/stretchr/testify/assert"
)

const shortTime = 50 * time.Millisecond

func TestBufferChanInput(t *testing.T) {
	c := New()
	c.Close() // We don't want the channel to process messages
	res := tryWithTimeout(shortTime, func () {
		c.Input() <- 1
	})
	assert.False(t, res)

	c = NewWithBuffer(1)
	c.Close()
	res = tryWithTimeout(shortTime, func () {
		c.Input() <- 1
	})
	assert.True(t, res)
	res = tryWithTimeout(shortTime, func () {
		c.Input() <- 1
	})
	assert.False(t, res)
}

func TestListenerCreationAndDeletion(t *testing.T) {
	c := New()
	assert.Equal(t, 0, syncMapLength(c.allListeners))

	lis1 := c.Listen()
	assert.Equal(t, 1, syncMapLength(c.allListeners))
	c.allListeners.Range(func (k, v interface{}) bool {
		assert.Equal(t, lis1, k)
		return true
	})

	lis2 := c.Listen()
	assert.Equal(t, 2, syncMapLength(c.allListeners))
	var last interface{} = nil
	c.allListeners.Range(func (k, v interface{}) bool {
		assert.True(t, (k == lis1 || k == lis2) && k != last)
		last = k
		return true
	})

	lis1.Close()
	assert.Equal(t, 1, syncMapLength(c.allListeners))
	c.allListeners.Range(func (k, v interface{}) bool {
		assert.Equal(t, lis2, k)
		return true
	})
}

func TestListenerMessages(t *testing.T) {
	c := New()

	type Msg struct {
		lis *Listener
		num int
	}
	
	lis1 := c.Listen()
	lis2 := c.Listen()
	outputChan := make(chan interface{}, 2)

	listen := func (l *Listener) {
		for {
			msg := <- l.Output()
			outputChan <- Msg{
				lis: l,
				num: msg.(int),
			}
		}
	}

	go listen(lis1)
	go listen(lis2)

	c.Input() <- 1
	assert.Equal(t, 1, (<-outputChan).(Msg).num)
	assert.Equal(t, 1, (<-outputChan).(Msg).num)

	c.Input() <- 2
	assert.Equal(t, 2, (<-outputChan).(Msg).num)
	assert.Equal(t, 2, (<-outputChan).(Msg).num)

	lis1.Close()

	c.Input() <- 3
	o := <-outputChan
	assert.Equal(t, 3, o.(Msg).num)
	assert.Equal(t, lis2, o.(Msg).lis)

	res := tryWithTimeout(shortTime, func() {
		<-outputChan
	})
	assert.False(t, res)
}

func TestUntilCloseListener(t *testing.T) {
	c := New()
	lis := c.Listen()

	out := make(chan interface{}, 1)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-lis.UntilClose():
				close(done)
				return
			case msg := <-lis.Output():
				out <- msg
			}
		}
	}()

	c.Input() <- 1
	assert.Equal(t, 1, (<-out).(int))

	c.Close()
	res := tryWithTimeout(shortTime, func () {
		<- done
	})

	assert.True(t, res)
}

func TestInfiniteChan(t *testing.T) {
	c := NewInfinite()
	lisA := c.Listen()
	lisB := c.Listen()

	res := tryWithTimeout(shortTime, func () {
		c.Input() <- 1
	})
	assert.True(t, res)
	res = tryWithTimeout(shortTime, func () {
		c.Input() <- 2
	})
	assert.True(t, res)
	res = tryWithTimeout(shortTime, func () {
		c.Input() <- 3
	})
	assert.True(t, res)

	res = tryWithTimeout(shortTime, func () {
		<- lisA.Output()
		<- lisA.Output()
	})
	assert.False(t, res)

	lisA.Close()
	lisB.Close()
}

func TestInfiniteListeners(t *testing.T) {
	c := New()
	lisA := c.ListenInfinite()
	lisB := c.ListenInfinite()

	res := tryWithTimeout(shortTime, func () {
		c.Input() <- 1
	})
	assert.True(t, res)
	res = tryWithTimeout(shortTime, func () {
		c.Input() <- 2
	})
	assert.True(t, res)
	res = tryWithTimeout(shortTime, func () {
		c.Input() <- 3
	})
	assert.True(t, res)

	res = tryWithTimeout(shortTime, func () {
		a := <- lisA.Output()
		b := <- lisA.Output()
		assert.Equal(t, 1, a)
		assert.Equal(t, 2, b)
	})
	assert.True(t, res)

	res = tryWithTimeout(shortTime, func () {
		a := <- lisB.Output()
		b := <- lisB.Output()
		assert.Equal(t, 1, a)
		assert.Equal(t, 2, b)
	})
	assert.True(t, res)

	lisA.Close()
	lisB.Close()
}