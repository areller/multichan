package main

import (
	"time"
	"os/signal"
	"os"
	"strconv"
	"fmt"
	"github.com/areller/multichan"
)

type Producer struct {
	Event *multichan.Chan
}

func NewProducer() *Producer {
	return &Producer{
		Event: multichan.New(),
	}
}

type Consumer struct {
	Event *multichan.Listener
}

func (c *Consumer) run() {
	for {
		select {
		case msg := <-c.Event.Output():
			fmt.Println(msg.(string))
		case <-c.Event.UntilClose():
			fmt.Println("Producer closing propagated to this consumer")
			return
		}
	}
}

func NewConsumer(producer *Producer) *Consumer {
	c := &Consumer{
		Event: producer.Event.Listen(),
	}

	go c.run()
	return c
}

func main() {
	p := NewProducer()
	NewConsumer(p)
	NewConsumer(p)

	c := 0
	go func() {
		for {
			p.Event.Input() <- "Message " + strconv.Itoa(c)
			c++
			time.Sleep(time.Second)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	<- sigChan
	p.Event.Close()
	time.Sleep(time.Second)
}