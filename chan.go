package multichan

import (
	"sync"
)

type Chan struct {
	closeChan chan int
	mu sync.RWMutex
	inputChan chan interface{}
	allListeners map[*Listener]bool
}

func (mc *Chan) sendToAll(msg interface{}) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	for l := range mc.allListeners {
		l.outputChan <- msg
	}
}

func (mc *Chan) closeAllListeners() {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	for l := range mc.allListeners {
		if l.withClose {
			close(l.closeChan)
		}
	}
}

func (mc *Chan) run() {
	r := true
	for r {
		select {
		case msg := <-mc.inputChan:
			mc.sendToAll(msg)
		case <-mc.closeChan:
			r = false	
		}
	}
	
	mc.closeAllListeners()
	mc.closeChan <- 0
}

func (mc *Chan) Close() {
	mc.closeChan <- 0
	<-mc.closeChan
}

func (mc *Chan) Listen(withClose bool) *Listener {
	lis := newListener(mc, withClose)

	mc.mu.Lock()
	mc.allListeners[lis] = true
	mc.mu.Unlock()

	return lis
}

func (mc *Chan) Input() chan interface{} {
	return mc.inputChan
}

func New() *Chan {
	return NewWithBuffer(0)
}

func NewWithBuffer(bufferSize int) *Chan {
	c := &Chan{
		inputChan: make(chan interface{}, bufferSize),
		allListeners: make(map[*Listener]bool),
		closeChan: make(chan int),
	}	

	go c.run()
	return c
}