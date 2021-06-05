package rwlocker

import (
	"sync"
	"testing"
)

func TestWriteLock(t *testing.T) {
	l := NewRWLocker()

	var c int
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Lock()
			c += 1
			l.Unlock()
		}()
	}

	wg.Wait()

	if c != 10 {
		t.Fail()
	}
}
