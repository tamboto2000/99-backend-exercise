package gocb

import (
	"sync"
	"time"
)

type errCounterStatus struct {
	mx *sync.Mutex
	t  bool
}

func newErrCounterStatus() *errCounterStatus {
	return &errCounterStatus{
		mx: new(sync.Mutex),
	}
}

func (ecs *errCounterStatus) timerOn(b bool) {
	defer ecs.mx.Unlock()
	ecs.mx.Lock()
	ecs.t = b
}

func (ecs *errCounterStatus) isOn() bool {
	defer ecs.mx.Unlock()
	ecs.mx.Lock()
	b := ecs.t
	return b
}

type errCounter struct {
	c        *counter
	interval time.Duration
	status   *errCounterStatus
	r        chan struct{}
}

func newErrCounter(interval time.Duration) *errCounter {
	ec := errCounter{
		c:        newCounter(),
		interval: interval,
		status:   newErrCounterStatus(),
		r:        make(chan struct{}),
	}

	return &ec
}

func (ec *errCounter) count() int64 {
	return ec.c.count()
}

func (ec *errCounter) incr() {
	ec.c.incr()
	if !ec.status.isOn() {
		go ec.startTimer()
	}
}

func (ec *errCounter) reset() {
	if ec.status.isOn() {
		ec.r <- struct{}{}
	}
}

func (ec *errCounter) startTimer() {
	ec.status.timerOn(true)
	defer func() {
		ec.c.reset()
		ec.status.timerOn(false)
	}()

	tr := time.NewTimer(ec.interval)
	select {
	case <-tr.C:
		return
	case <-ec.r:
		if !tr.Stop() {
			go func() {
				<-tr.C
			}()
		}

		return
	}
}
