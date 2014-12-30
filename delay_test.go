package delay

import (
	"testing"
	"time"
)

func TestNewDelayerStartsWithNoPendingTimers(t *testing.T) {
	d := NewDelayer(func(key, payload string) {}, 0)

	got := d.Pending()
	if got != 0 {
		t.Errorf("Expected 0 pending, got %d\n", got)
	}
}

func TestRegisterTriggersCallbackAfterSpecifiedTime(t *testing.T) {
	sem := make(chan bool)
	cb := func(key, payload string) {
		sem <- true
	}

	d := NewDelayer(cb, 10*time.Millisecond)

	d.Register("a", "message")
	start := time.Now()
	<-sem

	elapsed := time.Since(start)
	if elapsed <= 10*time.Millisecond {
		t.Errorf("Callback should run after 10ms, %s elasped\n", elapsed)
	}
}

func TestRegisteredCallbacksAreCanceled(t *testing.T) {
	d := NewDelayer(func(key, payload string) {}, 1*time.Second)
	d.Register("a", "message")

	d.Cancel("a")

	got := d.Pending()
	if got != 0 {
		t.Errorf("It should not be pending callbacks, got %d\n", got)
	}
}

func TestCallbacksAreUpdated(t *testing.T) {
	c := make(chan string)

	d := NewDelayer(func(key, payload string) {
		c <- payload
	}, 10*time.Millisecond)

	d.Register("a", "1")
	d.Register("a", "2")
	d.Register("a", "3")

	got := <-c
	if got != "3" {
		t.Errorf("Message expected '%s', got '%s'", "3", got)
	}
}

func TestFlushTriggerAllPendingCallbacks(t *testing.T) {
	pending := 2
	c := make(chan string, pending)

	d := NewDelayer(func(key, payload string) {
		c <- payload
	}, 10*time.Second)

	d.Register("a", "message")
	d.Register("b", "message")
	d.Flush()

loop:
	for {
		select {
		case <-time.After(10 * time.Millisecond):
			t.Errorf("Messages was not flushed")
			return
		case <-c:
			pending = pending - 1
			if pending == 0 {
				break loop
			}
		}
	}

	got := d.Pending()
	if got != pending {
		t.Errorf("got: %d, expected: %d\n", got, pending)
	}
}

func TestFlushTriggerASpecificCallback(t *testing.T) {
	c := make(chan string, 2)

	d := NewDelayer(func(key, payload string) {
		c <- payload
	}, 10*time.Second)

	d.Register("a", "1")
	d.Register("b", "2")

	d.Flush("a")

loop:
	for {
		select {
		case <-time.After(10 * time.Millisecond):
			break loop
		case <-c:
		}
	}

	got := d.Pending()
	if got != 1 {
		t.Errorf("Pending() got: %d, expected: 1\n", got)
	}
}
