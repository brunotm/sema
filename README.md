# Sema

Sema is a simple semaphore implementation for Go (golang) using channels to control concurrent access to shared resources.
(Go still lacks a user accessible semaphore implementation in the standard library)

## Example

```go

package main
import (
	"time"
	"github.com/brunotm/sema"
)

func main() {
	max := 100
	s := sema.New(max)

	s.Acquire()
	defer s.Release()
	// DO WORK...

	// OR
	if s.TryAcquire() {
		defer s.Release()
		// DO WORK...
	}

	// OR
	if s.AcquireWithin(5 * time.Millisecond) {
		defer s.Release()
		// DO WORK...
	}

}

```
Written by Bruno Moura <brunotm@gmail.com>