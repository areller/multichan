package multichan

type Listener struct {
	mainChan *Chan
	outputChan chan interface{}
	withClose bool
	closeChan chan struct{}
}

func (l *Listener) Close() {
	l.mainChan.mu.Lock()
	delete(l.mainChan.allListeners, l)
	l.mainChan.mu.Unlock()
}

func (l *Listener) UntilClose() <-chan struct{} {
	return l.closeChan
}

func (l *Listener) Output() <-chan interface{} {
	return l.outputChan
}

func newListener(mc *Chan, withClose bool) *Listener {
	return &Listener{
		mainChan: mc,
		outputChan: make(chan interface{}),
		withClose: withClose,
		closeChan: make(chan struct{}),
	}
}