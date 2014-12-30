package delay

import "time"

type Delayer struct {
	After    time.Duration
	Callback func(key, payload string)
	timers   map[string]*time.Timer
}

func NewDelayer(fn func(key, payload string), after time.Duration) *Delayer {
	delayer := &Delayer{After: after, Callback: fn}
	delayer.timers = make(map[string]*time.Timer)
	return delayer
}

func (d *Delayer) Register(key, payload string) {
	if timer, ok := d.timers[key]; ok {
		timer.Stop()
	}

	d.timers[key] = time.AfterFunc(d.After, func() {
		d.Callback(key, payload)
		delete(d.timers, key)
	})
}

func (d *Delayer) Cancel(key string) bool {
	timer, ok := d.timers[key]
	if !ok {
		return false
	}

	timer.Stop()
	delete(d.timers, key)
	return true
}

func (d *Delayer) Pending() int {
	return len(d.timers)
}

func (d *Delayer) Flush(keys ...string) int {
	var timers map[string]*time.Timer

	if len(keys) > 0 {
		timers = make(map[string]*time.Timer)
		for _, key := range keys {
			if timer, ok := d.timers[key]; ok {
				timers[key] = timer
			}
		}
	} else {
		// If no keys are specified, we flush all of them
		timers = d.timers
	}

	flushed := 0
	for _, timer := range timers {
		if timer.Reset(0) == true {
			flushed = flushed + 1
		}
	}

	return flushed
}