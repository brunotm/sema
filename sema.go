package sema

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	errCap = fmt.Errorf("sema: capacity must be at least 1")
)

// Sema provides a simple semaphore implementation using a channel
type Sema struct {
	sm      chan struct{}
	holders *int64
}

// New creates a new semaphore with the given maximum capacity for concurrent access
func New(cap int) (*Sema, error) {
	if cap < 1 {
		return nil, errCap
	}
	h := int64(0)
	return &Sema{
		sm:      make(chan struct{}, cap),
		holders: &h,
	}, nil
}

// Acquire the semaphore, will block if semaphore is full
// until any other holder release it
func (s *Sema) Acquire() {
	s.checkNil()
	s.sm <- struct{}{}
	atomic.AddInt64(s.holders, 1)
}

// Release the semaphore
func (s *Sema) Release() {
	s.checkNil()
	if atomic.AddInt64(s.holders, -1) < 0 {
		panic("sema: calling release on a empty semaphore")
	}
	<-s.sm
}

// TryAcquire the semaphore without blocking
// return true on success and false on failure
func (s *Sema) TryAcquire() bool {
	s.checkNil()
	select {
	case s.sm <- struct{}{}:
		atomic.AddInt64(s.holders, 1)
		return true
	default:
		return false
	}
}

// AcquireWithin the given timeout
// return true on success and false on failure
func (s *Sema) AcquireWithin(timeout time.Duration) bool {
	s.checkNil()
	select {
	case s.sm <- struct{}{}:
		atomic.AddInt64(s.holders, 1)
		return true
	case <-time.After(timeout):
		return false
	}
}

// Holders return the current holders count
func (s *Sema) Holders() int {
	s.checkNil()
	return int(atomic.LoadInt64(s.holders))
}

// Cap return semaphore capacity
func (s *Sema) Cap() int {
	s.checkNil()
	return cap(s.sm)
}

func (s *Sema) checkNil() {
	if s.sm == nil {
		panic("sema: calling on a nil semaphore")
	}
}
