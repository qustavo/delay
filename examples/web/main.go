package main

import (
	"fmt"
	"time"

	"github.com/gchaincl/delay"
	"github.com/gchaincl/delay/web"
)

func main() {

	delayer := delay.NewDelayer(func(key, payload string) {
		fmt.Printf("[%s] -> %s\n", key, payload)
	}, 2*time.Second)

	s := web.NewServer()
	s.Handle("emails", delayer)
	panic(s.Listen(":3000"))
}
