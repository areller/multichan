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

	close(l.closeChan)
}

func (l *Listener) UntilClose() <-chan struct{} {
	return l.closeChan
}

func (l *Listener) Output() <-chan interface{} {
	return l.outputChan
}

func (l *Listener) Attach(del func (msg interface{})) {
	go func(l *Listener, del func (msg interface{})) {
		for {
			select {
			case <-l.UntilClose():
				return
			case msg := <-l.Output():
				del(msg)
			}
		}
	}(l, del)
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