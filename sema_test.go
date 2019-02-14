package sema

import (
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	if _, err := New(0); err == nil {
		t.Errorf("New() error = %v, wantErr %v", err, errCap)
	}

	if _, err := New(1); err != nil {
		t.Errorf("New() error = %v, wantErr %v", err, nil)
	}
}

func TestAcquireRelease(t *testing.T) {
	sema, _ := New(10)

	for x := 0; x < sema.Cap(); x++ {
		sema.Acquire()
	}

	if sema.Holders() != sema.Cap() {
		t.Errorf("Holders() %d != %d Sema.Cap() ", sema.Holders(), sema.Cap())
	}

	for x := 1; x <= sema.Cap(); x++ {
		sema.Release()
	}

	if sema.Holders() != 0 {
		t.Errorf("Holders() == %d  after releasing Sema.Cap() == %d", sema.Holders(), sema.Cap())
	}

}

func TestTryAcquire(t *testing.T) {
	sema, _ := New(1)

	if !sema.TryAcquire() {
		t.Errorf("Failed to TryAcquire() with Cap() == %d and Holders() == %d", sema.Cap(), sema.Holders())
	}

	if sema.TryAcquire() {
		t.Errorf("Success to TryAcquire() with Cap() == %d and Holders() == %d", sema.Cap(), sema.Holders())
	}
}

func TestAcquireWithin(t *testing.T) {
	sema, _ := New(1)

	if !sema.AcquireWithin(1 * time.Millisecond) {
		t.Errorf("Failed to AcquireWithin() with Cap() == %d and Holders() == %d", sema.Cap(), sema.Holders())
	}
	if sema.AcquireWithin(1 * time.Millisecond) {
		t.Errorf("Success to AcquireWithin() with Cap() == %d and Holders() == %d", sema.Cap(), sema.Holders())
	}
}

func TestConcurrency(t *testing.T) {
	cap := 10
	sema, _ := New(cap)
	wg := &sync.WaitGroup{}

	for x := 1; x <= cap; x++ {
		sema.Acquire()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for x := 1; x <= cap; x++ {
			sema.Acquire()
			time.Sleep(time.Duration(cap) * time.Nanosecond)
			sema.Release()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for x := 1; x <= cap; x++ {
			sema.Acquire()
			time.Sleep(time.Duration(cap) * time.Nanosecond)
			sema.Release()
		}
	}()

	for x := 1; x <= cap; x++ {
		time.Sleep(time.Duration(cap) * time.Nanosecond)
		sema.Release()
	}

	wg.Wait()
	if sema.Holders() != 0 {
		t.Errorf("Expected Holders() == %d, got %d", 0, sema.Holders())
	}
}

func BenchmarkAcquireRelease(b *testing.B) {
	sema, _ := New(10)
	for i := 0; i < b.N; i++ {
		for x := 0; x < 10; x++ {
			sema.Acquire()
		}

		for x := 0; x < 10; x++ {
			sema.Release()
		}
	}
}

func BenchmarkConcurrency(b *testing.B) {
	cap := 10
	sema, _ := New(cap)
	wg := &sync.WaitGroup{}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for x := 1; x <= cap; x++ {
			sema.Acquire()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for x := 1; x <= cap; x++ {
				sema.Acquire()
				sema.Release()
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for x := 1; x <= cap; x++ {
				sema.Acquire()
				sema.Release()
			}
		}()

		for x := 1; x <= cap; x++ {
			sema.Release()
		}

		wg.Wait()
	}
}
