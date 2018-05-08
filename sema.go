package sema

import (
	"fmt"
	"time"
)

var (
	errCap = fmt.Errorf("sema: capacity must be at least 1")
)

// Sema provides a simple semaphore implementation using a channel
type Sema chan struct{}

// New creates a new semaphore with the given maximum capacity for concurrent access.
// Will panic if cap < 1
func New(cap int) (Sema, error) {
	if cap < 1 {
		return nil, errCap
	}
	return make(Sema, cap), nil
}

// Acquire the semaphore, will block if semaphore is full
// until any other holder release it.
// Will panic if semaphore is nil
func (s Sema) Acquire() {
	s.checkNil()
	s <- struct{}{}
}

// Release the semaphore
// Will panic if called on a non acquired semaphore
func (s Sema) Release() {
	s.checkNil()
	if len(s) < 1 {
		panic("sema: calling release on a empty semaphore")
	}
	<-s
}

// TryAcquire the semaphore without blocking return true on success and false on failure.
// Will panic if semaphore is nil
func (s Sema) TryAcquire() bool {
	s.checkNil()
	select {
	case s <- struct{}{}:
		return true
	default:
		return false
	}
}

// AcquireWithin the given timeout return true on success and false on failure
// Will panic if semaphore is nil
func (s Sema) AcquireWithin(timeout time.Duration) bool {
	s.checkNil()
	select {
	case s <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Holders return the current holders count
// Will panic if semaphore is nil
func (s Sema) Holders() int {
	s.checkNil()
	return len(s)
}

// Cap return semaphore capacity
// Will panic if semaphore is nil
func (s Sema) Cap() int {
	s.checkNil()
	return cap(s)
}

func (s Sema) checkNil() {
	if s == nil {
		panic("sema: calling on a nil semaphore")
	}
}
