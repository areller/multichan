# MultiChan
Simple Library for Multi-Listener Channels in Go.
Imagine a regular go channel but one that can have many consumer, all receiving the same messages.

## Why
This is a pattern that i find myself using a lot.
This library is supposed to get rid of all the boilerplate, in my projects and hopefully in yours too.

## Producer
```golang
c := multichan.New()
```

## Consumer(s)
```golang
lis1 := c.Listen()
lis2 := c.Listen()
listen := func (l *multichan.Listener) {
    for {
        msg := <- l.Output()
        fmt.Println(msg.(string))
    }
}

go listen(lis1)
go listen(lis2)
```

## Producing Messages
```golang
c.Input() <- "Message a"
c.Input() <- "Message b"
```

## Output
```
Message a
Message a
Message b
Message b
```

## Closing a Listener
If a thread is no longer listening, you may want to stop its listener by calling the `Close` method.

```golang
listener := c.Listen()
// Do some listening
// ...
listener.Close()
```

This will inform the channel (`c`) to not queue messages to that listener anymore.

## `UntilClose`
When a channel closes, all of its listeners close subsequently.
You can subscribe to `listener.UntilClose()` to receive a notification for that.

```golang
c := multichan.New()
listen := func (num string) {
    lis := c.Listen()
    <-lis.UntilClose()
    fmt.Println("Listener #" + num + " was closed due to its channel closing")
}
listen("1")
listen("2")

// Do something with channels, listeners
// ...

c.Close()
```

The output would be

```
Listener #1 was closed due to its channel closing
Listener #2 was closed due to its channel closing
```

## Buffering
You can create buffered channels and listeners
```golang
c := multichan.NewWithBuffer(2)
```
```golang
lis1 := c.Listen()
lis2 := c.Listen()
```

In this example the channel is buffered (buffer of 2) but the listeners are not meaning that you have to pull from all listeners simultaneously if you want to see all messages

```golang
c.Input() <- 1
c.Input() <- 2 // Won't block since channel is buffered

<-lis1.Output() // Might block
<-lis1.Output() // Will definitely block
```

Messages that are queued to the channel are queued to the listeners in sequence, if `lisA` happens to be the first in the sequence of listeners, the first pull from its Output() won't block, otherwise it will.

If you want all listeners to receive all messages independently regardless of whether all listeners are active or not, **all** listeners will have to be buffered.

```golang
lis1 := c.ListenBuffered(2)
lis2 := c.ListenBuffered(2)
```

Then,

```golang
c.Input() <- 1
c.Input() <- 2 // Won't block since channel is buffered

<-lis1.Output() // Won't block
<-lis1.Output() // Won't block

<-lis2.Output() // Won't block
<-lis2.Output() // Won't block
```

**Note** If just 1 listener out of 100 is not buffered, it will block the rest from receiving the latest messages unless it is being polled constantly.

## Infinite Channels
Multichan support infinite channels (i.e. channels with infinite buffer)
```golang
c := multichan.NewInfinite()
```
```golang
listener := c.ListenInfinite()
```

## More Examples
Go to [this](./examples) folder for more examples