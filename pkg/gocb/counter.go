package gocb

import (
	"sync/atomic"
)

type counter struct {
	*atomic.Int64
}

func newCounter() *counter {
	return &counter{new(atomic.Int64)}
}

func (c *counter) count() int64 {
	return c.Load()
}

func (c *counter) incr() {
	c.Add(1)
}

func (c *counter) reset() {
	c.Store(0)
}
