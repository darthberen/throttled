// Package throttled implements various helpers to manage the lifecycle of goroutines.
package throttled

// WaitGroup limits the number of concurrent goroutines.
type WaitGroup struct {
	throttle    int
	completed   chan bool
	outstanding int
}

// NewWaitGroup instantiates a new WaitGroup with the given throttle.
func NewWaitGroup(throttle int) *WaitGroup {
	return &WaitGroup{
		outstanding: 0,
		throttle:    throttle,
		completed:   make(chan bool, throttle),
	}
}

// Add will block until the number of goroutines being throttled
// has fallen below the throttle
func (w *WaitGroup) Add() {
	if w.outstanding+1 > w.throttle {
		select {
		case <-w.completed:
			w.outstanding--
			return
		}
	}
	w.outstanding++
}

// Done signal that a goroutine has completed
func (w *WaitGroup) Done() {
	w.completed <- true
}

// Wait until all of the throttled goroutines have signaled they are done
func (w *WaitGroup) Wait() {
	if w.outstanding == 0 {
		return
	}
	for w.outstanding > 0 {
		select {
		case <-w.completed:
			w.outstanding--
		}
	}
}
