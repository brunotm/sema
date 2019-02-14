package sema

import (
	"errors"
	"sync"
	"time"
)

var (
	errCap = errors.New("sema: capacity must be at least 1")
)

// Sema provides a simple semaphore implementation
type Sema struct {
	cond  *sync.Cond
	cap   int
	count int
}

// New creates a new semaphore with the given maximum capacity for concurrent access.
func New(cap int) (s *Sema, err error) {
	if cap < 1 {
		return nil, errCap
	}
	s = &Sema{}
	s.cap = cap
	s.cond = sync.NewCond(&sync.Mutex{})
	return s, nil
}

// Acquire the semaphore, will block if semaphore is full until any other holder release it.
func (s *Sema) Acquire() {
	s.cond.L.Lock()
	for s.count == s.cap {
		s.cond.Wait()
	}
	s.count++
	s.cond.Signal()
	s.cond.L.Unlock()
}

// Release the semaphore allowing waking waiters if any to acquire.
func (s *Sema) Release() {
	s.cond.L.Lock()
	if s.count == 0 {
		panic("sema: calling release on a empty semaphore")
	}
	s.count--
	s.cond.Signal()
	s.cond.L.Unlock()
}

// TryAcquire the semaphore without blocking return true on success and false on failure.
func (s *Sema) TryAcquire() bool {
	s.cond.L.Lock()
	if s.count == s.cap {
		s.cond.L.Unlock()
		return false
	}
	s.count++
	s.cond.Signal()
	s.cond.L.Unlock()
	return true
}

// AcquireWithin the given timeout return true on success and false on failure
func (s *Sema) AcquireWithin(timeout time.Duration) bool {
	if s.TryAcquire() {
		return true
	}
	time.Sleep(timeout)
	return s.TryAcquire()
}

// Holders return the current holders count
func (s *Sema) Holders() int {
	return s.count
}

// Cap return semaphore capacity
func (s *Sema) Cap() int {
	return s.cap
}
