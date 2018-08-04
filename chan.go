package multichan

import (
	"sync"
)

type Closable interface {
	Close()
}

type Chan struct {
	closeChan chan int
	inputChan chan interface{}
	outputChan chan interface{}
	allListeners *sync.Map
	closable Closable
}

func (mc *Chan) sendToAll(msg interface{}) {
	mc.allListeners.Range(func (k, v interface{}) bool {
		k.(*Listener).inputChan <- msg
		return true
	})
}

func (mc *Chan) closeAllListeners() {
	mc.allListeners.Range(func (k, v interface{}) bool {
		close(k.(*Listener).closeChan)
		return true
	})
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
	mc.allListeners.Store(lis, true)

	return lis
}

func (mc *Chan) ListenBuffered(bufferSize int) *Listener {
	lis := newBufferedListener(mc, bufferSize)
	mc.allListeners.Store(lis, true)

	return lis
}

func (mc *Chan) ListenInfinite() *Listener {
	lis := newInfiniteListener(mc)
	mc.allListeners.Store(lis, true)

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
		allListeners: new(sync.Map),
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
		allListeners: new(sync.Map),
		closeChan: make(chan int),
		closable: inf,
	}

	go c.run()
	return c
}