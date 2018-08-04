package multichan

type infiniteChan struct {
	inChan chan interface{}
	outChan chan interface{}
	closeChan chan struct{}
	queue []interface{}
}

func (ic *infiniteChan) out() chan interface{} {
	if len(ic.queue) == 0 {
		return nil
	}
	return ic.outChan
}

func (ic *infiniteChan) dequeue() interface{} {
	if len(ic.queue) == 0 {
		return nil
	}
	return ic.queue[0]
}

func (ic *infiniteChan) cut() {
	if len(ic.queue) != 0 {
		ic.queue = ic.queue[1:]
	}
}

func (ic *infiniteChan) run() {
	for {
		select {
		case <-ic.closeChan:
			return
		case obj := <-ic.inChan:
			ic.queue = append(ic.queue, obj)
		case ic.out() <-ic.dequeue():
			ic.cut()
		}
	}
}

func (ic *infiniteChan) Close() {
	close(ic.closeChan)
}

func newInfinite() *infiniteChan {
	ic := &infiniteChan{
		inChan: make(chan interface{}),
		outChan: make(chan interface{}),
		closeChan: make(chan struct{}, 1),
	}

	go ic.run()
	return ic
}