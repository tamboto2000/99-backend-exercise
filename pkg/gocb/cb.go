package gocb

import (
	"errors"
	"time"
)

type CircuitStatus uint

const (
	Closed CircuitStatus = iota
	Open
	HalfOpen
)

var ErrCircuitOpen = errors.New("circuit is open")

// Circuit provides APIs for circuit breaker implementation
type Circuit interface {
	// Name return the name of the circuit breaker
	Name() string
	// Status return current status of the circuit
	Status() CircuitStatus
	// Run run the function or routine that you want to execute
	Run(func() error) error
}

// Settings contains settings to configure
// CircuitBreaker
type Settings struct {
	// ErrThreshold sets threshold for error rate.
	// Once error rate is equal to this value, circuit will trip
	// which change the circuit status to Open
	ErrThreshold uint
	// ErrInterval sets how much time should passed before error counter
	// reset after first error detection.
	// For instance, should you set ErrInterval to 10 seconds, and
	// there is another error before 10 seconds after the first error, error counter
	// will detect it and increment the error count by 1, giving 2 errors in total,
	// which if your ErrThreshold is 2, your
	// circuit will trip
	ErrInterval time.Duration
	// Timeout sets the duration of how long the circuit will be in
	// Open state before it can switch to HalfOpen
	Timeout time.Duration
	// Retry sets how much successful retries before circuit state
	// can be switched from HalfOpen to Closed
	Retry uint
	// OnClosed sets callback for when circuit state is switched to Closed
	OnClosed func()
	// OnOpen sets callback for when circuit state is switched to Open
	OnOpen func()
	// OnHalfOpen sets callback for when circuit state is switched to HalfOpen
	OnHalfOpen func()
}

// CircuitBreaker implements Circuit Breaker patterns.
type CircuitBreaker struct {
	name         string
	settings     Settings
	errCounter   *errCounter
	retryCounter *counter
	statStore    *statusStore
}

// NewCircuitBreaker initiate new CircuitBreaker.
// Use settings to configure it, you can also give it a name.
func NewCircuitBreaker(name string, settings Settings) *CircuitBreaker {
	statStore := newStatusStore()
	errCounter := newErrCounter(settings.ErrInterval)

	statStore.set(Closed)
	cb := CircuitBreaker{
		name:         name,
		settings:     settings,
		errCounter:   errCounter,
		retryCounter: newCounter(),
		statStore:    statStore,
	}

	return &cb
}

// Name return the name of the circuit breaker
func (cb *CircuitBreaker) Name() string {
	return cb.name
}

// Status return current circuit status
func (cb *CircuitBreaker) Status() CircuitStatus {
	return cb.statStore.status()
}

// Run run a function or routine
func (cb *CircuitBreaker) Run(f func() error) error {
	return cb.run(f)
}

func (cb *CircuitBreaker) run(f func() error) error {
	if cb.statStore.status() == Closed {
		err := f()
		if err != nil {
			cb.errCounter.incr()
			if cb.errCounter.count() == int64(cb.settings.ErrThreshold) {
				cb.trip()
			}

			return err
		}

		cb.errCounter.reset()
		return nil
	}

	if cb.statStore.status() == Open {
		return ErrCircuitOpen
	}

	if cb.statStore.status() == HalfOpen {
		err := f()
		if err != nil {
			cb.trip()
			return err
		}

		cb.retryCounter.incr()
		if cb.retryCounter.count() == int64(cb.settings.Retry) {
			cb.restore()
		}
	}

	return nil
}

// trip sets circuit status to open
func (cb *CircuitBreaker) trip() {
	cb.statStore.set(Open)
	cb.reset()
	go cb.waitTimeout()

	if cb.settings.OnOpen != nil {
		go cb.settings.OnOpen()
	}
}

// waitTimeout waits until timeout passed which after that
// will update the circuit status to half-open
func (cb *CircuitBreaker) waitTimeout() {
	tr := time.NewTimer(cb.settings.Timeout)
	<-tr.C
	cb.statStore.set(HalfOpen)
	if cb.settings.OnHalfOpen != nil {
		go cb.settings.OnHalfOpen()
	}
}

// restore switch back the circuit to closed
func (cb *CircuitBreaker) restore() {
	cb.statStore.set(Closed)
	cb.reset()
	if cb.settings.OnClosed != nil {
		go cb.settings.OnClosed()
	}
}

func (cb *CircuitBreaker) reset() {
	cb.errCounter.reset()
	cb.retryCounter.reset()
}
