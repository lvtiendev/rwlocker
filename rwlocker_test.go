package rwlocker

import (
	"sync"
	"testing"
	"time"
)

func TestWriteLock(t *testing.T) {
	l := NewRWLocker()

	var c int
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(10 * time.Millisecond)
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

func TestMakeSureRUnlockWithoutLockDoNotStarveLaterWrite(t *testing.T) {
	ch := make(chan bool, 1)
	l := NewRWLocker()

	// pretend that there's a bug in code
	l.RUnlock()
	l.RUnlock()

	l.RLock()

	go func() {
		l.Lock()
		ch <- true
		l.Unlock()
	}()

	l.RUnlock()

	if b := <-ch; !b {
		t.Fail()
	}
}
