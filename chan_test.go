package multichan

import (
	"time"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBufferChanInput(t *testing.T) {
	c := New()
	c.Close() // We don't want the channel to process messages
	res := tryWithTimeout(time.Second, func () {
		c.Input() <- 1
	})
	assert.False(t, res)

	c = NewWithBuffer(1)
	c.Close()
	res = tryWithTimeout(time.Second, func () {
		c.Input() <- 1
	})
	assert.True(t, res)
	res = tryWithTimeout(time.Second, func () {
		c.Input() <- 1
	})
	assert.False(t, res)
}

func TestListenerCreationAndDeletion(t *testing.T) {
	
	c := New()
	assert.Len(t, c.allListeners, 0)

	lis1 := c.Listen(false)
	assert.Len(t, c.allListeners, 1)
	for k := range c.allListeners {
		assert.Equal(t, lis1, k)
	}

	lis2 := c.Listen(false)
	assert.Len(t, c.allListeners, 2)
	n := 0
	for k := range c.allListeners {
		switch n {
			case 0:
				assert.Equal(t, lis1, k)
			case 1:
				assert.Equal(t, lis2, k)
		}
		n++
	}

	lis1.Close()
	assert.Len(t, c.allListeners, 1)
	for k := range c.allListeners {
		assert.Equal(t, lis2, k)
	}
}

func TestListenerMessages(t *testing.T) {
	c := New()

	type Msg struct {
		lis *listener
		num int
	}
	
	lis1 := c.Listen(false)
	lis2 := c.Listen(false)
	outputChan := make(chan interface{}, 2)

	listen := func (l *listener) {
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

	res := tryWithTimeout(time.Second, func() {
		<-outputChan
	})
	assert.False(t, res)
}

func TestUntilCloseListener(t *testing.T) {
	c := New()
	lis := c.Listen(true)

	out := make(chan interface{}, 1)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-lis.closeChan:
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
	res := tryWithTimeout(time.Second, func () {
		<- done
	})

	assert.True(t, res)
}