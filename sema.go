package sema

import (
	"errors"
	"time"
)

var (
	errCap = errors.New("sema: capacity must be at least 1")
)

// Sema provides a simple semaphore implementation
type Sema chan struct{}

// New creates a new semaphore with the given maximum capacity for concurrent access.
func New(size int) (s Sema, err error) {
	if size < 1 {
		return nil, errCap
	}
	return make(Sema, size), nil
}

// Acquire will block if semaphore is full until any other holder release it.
func (s Sema) Acquire() {
	s.check()
	s <- struct{}{}
}

// Release the semaphore allowing wating waiters to acquire.
func (s Sema) Release() {
	s.check()
	if len(s) < 1 {
		panic("sema: calling release on a empty semaphore")
	}
	<-s
}

// TryAcquire the semaphore without blocking return true on success and false on failure.
func (s Sema) TryAcquire() (ok bool) {
	s.check()
	select {
	case s <- struct{}{}:
		return true
	default:
		return false
	}
}

// AcquireWithin the given timeout return true on success and false on failure
func (s Sema) AcquireWithin(timeout time.Duration) (ok bool) {
	s.check()
	select {
	case s <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Holders return the current holders count
func (s Sema) Holders() (count int) {
	return len(s)
}

// Cap return semaphore capacity
func (s Sema) Cap() (sinze int) {
	return cap(s)
}

func (s Sema) check() {
	if s == nil {
		panic("sema: calling on a nil semaphore")
	}
}
