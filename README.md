delay [![Build Status](https://travis-ci.org/gchaincl/delay.svg)](https://travis-ci.org/gchaincl/delay)
=====

Delay goroutines at will

Rationale
---
You're notifying events, say _follow_ on users, so followed users will get
notified when someone start following them, now imagine _someone_ follows
_someother_ and, immediately after, it unfollows, you don't want to send the
notification to the supposedly followed user.
`Delay` tries to solve that problem.

It allows to *delay* event triggering defining a waiting time.
While an event is waiting to be triggered it can be updated or even cancelled
(so it will never triggered).

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
