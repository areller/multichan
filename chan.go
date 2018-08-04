package multichan

import (
	"sync"
)

type Closable interface {
	Close()
}

type Chan struct {
	closeChan chan int
	mu sync.RWMutex
	inputChan chan interface{}
	outputChan chan interface{}
	allListeners map[*Listener]bool
	closable Closable
}

func (mc *Chan) sendToAll(msg interface{}) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	for l := range mc.allListeners {
		l.inputChan <- msg
	}
}

func (mc *Chan) closeAllListeners() {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	for l := range mc.allListeners {
		close(l.closeChan)
	}
}

func (mc *Chan) run() {
	r := true
	for r {
		select {
		case msg := <-mc.outputChan:
			mc.sendToAll(msg)
		case <-mc.closeChan:
			r = false	
		}
	}
	
	mc.closeAllListeners()
	if mc.closable != nil {
		mc.closable.Close()
	}

	mc.closeChan <- 0
}

func (mc *Chan) Close() {
	mc.closeChan <- 0
	<-mc.closeChan
}

func (mc *Chan) Listen() *Listener {
	lis := newListener(mc)

	mc.mu.Lock()
	mc.allListeners[lis] = true
	mc.mu.Unlock()

	return lis
}

func (mc *Chan) ListenInfinite() *Listener {
	lis := newInfiniteListener(mc)

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
	inAndOut := make(chan interface{}, bufferSize)
	c := &Chan{
		inputChan: inAndOut,
		outputChan: inAndOut,
		allListeners: make(map[*Listener]bool),
		closeChan: make(chan int),
	}	

	go c.run()
	return c
}

func NewInfinite() *Chan {
	inf := newInfinite()
	c := &Chan{
		inputChan: inf.inChan,
		outputChan: inf.outChan,
		allListeners: make(map[*Listener]bool),
		closeChan: make(chan int),
		closable: inf,
	}

	go c.run()
	return c
}