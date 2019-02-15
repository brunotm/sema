package sema

import (
	"context"
	"runtime"
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

func TestAcquireWith(t *testing.T) {
	sema, _ := New(1)

	ctx, cancel := context.WithCancel(context.Background())

	if !sema.AcquireWith(ctx) {
		t.Errorf("Failed to AcquireWith() with Cap() == %d and Holders() == %d and a valid context", sema.Cap(), sema.Holders())
	}

	cancel()
	if sema.AcquireWith(ctx) {
		t.Errorf("Success to AcquireWith() with Cap() == %d and Holders() == %d and a canceled context", sema.Cap(), sema.Holders())
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
	size := 10
	sema, _ := New(size)
	wg := &sync.WaitGroup{}

	for x := 0; x < size; x++ {
		sema.Acquire()
	}

	for x := 0; x < runtime.NumCPU(); x++ {
		wg.Add(1)
		go func() {
			for x := 0; x < size; x++ {
				sema.Acquire()
				time.Sleep(time.Duration(size) * time.Nanosecond)
				sema.Release()
			}
			wg.Done()
		}()
	}

	for x := 0; x < size; x++ {
		time.Sleep(time.Duration(size) * time.Nanosecond)
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
	size := 10
	sema, _ := New(size)
	wg := &sync.WaitGroup{}

	cpus := runtime.NumCPU()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for x := 0; x < size; x++ {
			sema.Acquire()
		}

		for x := 0; x < cpus; x++ {
			wg.Add(1)
			go func() {
				for x := 0; x < size; x++ {
					sema.Acquire()
					sema.Release()
				}
				wg.Done()
			}()
		}

		for x := 0; x < size; x++ {
			sema.Release()
		}

		wg.Wait()
	}
}
