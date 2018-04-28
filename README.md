# MultiChan
Simple Library for Multi-Listener Channels in Go.
Imagine a regular go channel but one that can have many consumer, all receiving the same messages.

## Producer
```golang
c := multichan.New()
```

## Consumer(s)
```golang
lis1 := c.Listen(false)
lis2 := c.Listen(false)
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

## More Examples
Go to /examples folder for more examples