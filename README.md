delay
=====

Delay goroutines at will

Example
=====

```go
func main() {
	delayer := NewDelayer(func(key, payload string) {
		println(key, payload)
	}, 2*time.Second)

	delayer.Register("a", "Msg A")
	delayer.Register("b", "Msg B")
	delayer.Register("c", "Msg C")
	delayer.Register("c", "Msg C1")
	delayer.Register("c", "Msg C2")
	delayer.Register("c", "Msg C3")

	for {
		println(delayer.Pending(), " currently in queue")
		time.Sleep(1 * time.Second)
	}

}
```
