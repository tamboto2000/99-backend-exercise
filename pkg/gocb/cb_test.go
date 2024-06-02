package gocb

import (
	"errors"
	"testing"
	"time"
)

var (
	milisec500 = 500 * time.Millisecond
	errDummy   = errors.New("error dummy")
)

var statusStrMap = map[CircuitStatus]string{
	Open:     "Open",
	HalfOpen: "HalfOpen",
	Closed:   "Closed",
}

func cbStatusStr(s CircuitStatus) string {
	return statusStrMap[s]
}

// MockTask is used to test the circuit breaker implementation
// logic
type MockTask struct {
	t *testing.T
	c chan CircuitStatus
}

func newMockTask(t *testing.T) *MockTask {
	return &MockTask{
		t: t,
		c: make(chan CircuitStatus),
	}
}

func (mt *MockTask) onOpen() {
	mt.c <- Open
}

func (mt *MockTask) onHalfOpen() {
	mt.c <- HalfOpen
}

func (mt *MockTask) onClosed() {
	mt.c <- Closed
}

func (mt *MockTask) waitExpectStatus(timeout time.Duration, status CircuitStatus) {
	tr := time.NewTimer(timeout)
	select {
	case <-tr.C:
		mt.t.Errorf("callback for status %s is never called", cbStatusStr(status))
		return
	case s := <-mt.c:
		if !tr.Stop() {
			go func() {
				<-tr.C
			}()
		}

		if s != status {
			mt.t.Errorf("got status %s, want status %s", cbStatusStr(s), cbStatusStr(status))
		}

		return
	}
}

func (mt *MockTask) clean() {
	close(mt.c)
}

func TestCircuitBreaker_Run(t *testing.T) {
	cb := NewCircuitBreaker("test-cb", Settings{
		ErrThreshold: 1,
		ErrInterval:  1 * time.Second,
		Timeout:      1 * time.Second,
		Retry:        1,
	})

	err := cb.Run(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("got err = %v, want err = nil", err)
	}
}

func TestCircuitBreaker_Run_fromClosedToOpen(t *testing.T) {
	errThreshold := 3
	errInterval := 2 * time.Second
	mt := newMockTask(t)
	cb := NewCircuitBreaker("test-cb", Settings{
		ErrThreshold: uint(errThreshold),
		ErrInterval:  errInterval,
		Timeout:      1 * time.Second,
		Retry:        1,
		OnOpen:       mt.onOpen,
	})

	defer mt.clean()

	// try to produce error without tripping the circuit
	for i := 0; i < errThreshold-1; i++ {
		err := cb.Run(func() error {
			return errDummy
		})

		if err != errDummy {
			t.Errorf("got err = %s, want err = %v", err, errDummy)
		}
	}

	if cb.Status() != Closed {
		t.Errorf("got status = %v, want status = %v", cbStatusStr(cb.Status()), cbStatusStr(Closed))
	}

	// try to trip the circuit by producing 1 more error
	err := cb.Run(func() error {
		return errDummy
	})

	if err == nil {
		t.Errorf("got err = nil, want err = %v", errDummy)
	}

	if err != nil {
		if err != errDummy {
			t.Errorf("got err = %v, want err = %v", err, errDummy)
		}
	}

	// on running 1 more routine, the error returned should be ErrCircuitOpen
	err = cb.Run(func() error { return nil })
	if err == nil {
		t.Errorf("got err = nil, want err = %v", ErrCircuitOpen)
	}

	if err != nil {
		if err != ErrCircuitOpen {
			t.Errorf("got err = %v, want err = %v", err, ErrCircuitOpen)
		}
	}

	mt.waitExpectStatus(errInterval+milisec500, Open)
}

func TestCircuitBreaker_Run_fromOpenToHalfOpenToClosed(t *testing.T) {
	errThreshold := 3
	errInterval := 2 * time.Second
	timeout := 1 * time.Second
	retry := 3
	mt := newMockTask(t)
	cb := NewCircuitBreaker("test-cb", Settings{
		ErrThreshold: uint(errThreshold),
		ErrInterval:  errInterval,
		Timeout:      timeout,
		Retry:        3,
		OnOpen:       mt.onOpen,
		OnHalfOpen:   mt.onHalfOpen,
		OnClosed:     mt.onClosed,
	})

	defer mt.clean()

	// try to trip the circuit
	for i := 0; i < errThreshold; i++ {
		err := cb.Run(func() error {
			return errDummy
		})

		if err == nil {
			t.Errorf("got err = nil, want err = %v", errDummy)
		}
	}

	if cb.Status() != Open {
		t.Errorf("got status = %v, want status = %v", cbStatusStr(cb.Status()), Open)
	}

	mt.waitExpectStatus(milisec500, Open)
	mt.waitExpectStatus(timeout+milisec500, HalfOpen)

	// try to produce consecutive success run but not enough to
	// close the circuit
	for i := 0; i < retry-1; i++ {
		err := cb.Run(func() error { return nil })
		if err != nil {
			t.Errorf("got err = %v, want err = nil", err)
		}
	}

	if cb.Status() != HalfOpen {
		t.Errorf("got status = %v, want status = %v", cbStatusStr(cb.Status()), HalfOpen)
	}

	// try to restore the circuit by running 1 more successful run
	err := cb.Run(func() error { return nil })
	if err != nil {
		t.Errorf("got err = %v, want err = nil", err)
	}

	if cb.Status() != Closed {
		t.Errorf("got status = %v, want status = %v", cbStatusStr(cb.Status()), Closed)
	}

	mt.waitExpectStatus(milisec500, Closed)
}

func TestCircuitBreaker_Run_fromHalfOpenToOpen(t *testing.T) {
	errThreshold := 3
	errInterval := 2 * time.Second
	timeout := 1 * time.Second
	retry := 3
	mt := newMockTask(t)
	cb := NewCircuitBreaker("test-cb", Settings{
		ErrThreshold: uint(errThreshold),
		ErrInterval:  errInterval,
		Timeout:      timeout,
		Retry:        3,
		OnOpen:       mt.onOpen,
		OnHalfOpen:   mt.onHalfOpen,
		OnClosed:     mt.onClosed,
	})

	defer mt.clean()

	// try to trip the circuit
	for i := 0; i < errThreshold; i++ {
		err := cb.Run(func() error {
			return errDummy
		})

		if err == nil {
			t.Errorf("got err = nil, want err = %v", errDummy)
		}
	}

	if cb.Status() != Open {
		t.Errorf("got status = %v, want status = %v", cbStatusStr(cb.Status()), Open)
	}

	mt.waitExpectStatus(milisec500, Open)
	mt.waitExpectStatus(timeout+milisec500, HalfOpen)

	// try to produce consecutive success run but not enough to
	// close the circuit
	for i := 0; i < retry-1; i++ {
		err := cb.Run(func() error { return nil })
		if err != nil {
			t.Errorf("got err = %v, want err = nil", err)
		}
	}

	// trip the circuit again by running unsuccessful run
	err := cb.Run(func() error { return errDummy })
	if err == nil {
		t.Errorf("got err = nil, want err = %v", errDummy)
	}

	mt.waitExpectStatus(milisec500, Open)
}

func TestCircuitBreaker_Name(t *testing.T) {
	name := "test-cb"
	cb := NewCircuitBreaker(name, Settings{})
	if cb.Name() != name {
		t.Errorf("got name = %v, want name = %v", cb.Name(), name)
	}
}
