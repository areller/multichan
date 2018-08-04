package multichan

type Listener struct {
	mainChan *Chan
	inputChan chan interface{}
	outputChan chan interface{}
	closeChan chan struct{}
	closable Closable
}

func (l *Listener) Close() {
	l.mainChan.allListeners.Delete(l)
	
	if l.closable != nil {
		l.closable.Close()
	}
}

func (l *Listener) UntilClose() <-chan struct{} {
	return l.closeChan
}

func (l *Listener) Output() <-chan interface{} {
	return l.outputChan
}

func newListener(mc *Chan) *Listener {
	return newBufferedListener(mc, 0)
}

func newBufferedListener(mc *Chan, bufferSize int) *Listener {
	c := make(chan interface{}, bufferSize)
	return &Listener{
		mainChan: mc,
		inputChan: c,
		outputChan: c,
		closeChan: make(chan struct{}, 1),
	}
}

func newInfiniteListener(mc *Chan) *Listener {
	inf := newInfinite()
	return &Listener{
		mainChan: mc,
		inputChan: inf.inChan,
		outputChan: inf.outChan,
		closeChan: make(chan struct{}, 1),
		closable: inf,
	}
}