package main

import (
	"os"
	"fmt"
	"time"
	"github.com/areller/multichan"
)

func main() {
	c := multichan.New()
	listener := c.ListenInfinite()

	fastTicker := time.NewTicker(10 * time.Millisecond)
	slowTicker := time.NewTicker(time.Second)
	done := time.After(10 * time.Second)

	for {
		select {
		case t := <-fastTicker.C:
			c.Input() <- t
		case <-slowTicker.C:
			{
				t := (<-listener.Output()).(time.Time)
				fmt.Println("Pulled, was inserted at " + t.String())
			}
		case <-done:
			os.Exit(0)
		}
	}
}