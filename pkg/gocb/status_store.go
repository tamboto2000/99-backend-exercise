package gocb

import "sync/atomic"

type statusStore struct {
	*atomic.Value
}

func newStatusStore() *statusStore {
	return &statusStore{new(atomic.Value)}
}

func (ss *statusStore) set(status CircuitStatus) {
	ss.Store(status)
}

func (ss *statusStore) status() CircuitStatus {
	return ss.Load().(CircuitStatus)
}
