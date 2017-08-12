package sema

import (
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

	for x := 1; x <= sema.Cap(); x++ {
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
