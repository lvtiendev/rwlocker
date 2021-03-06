package rwlocker

import "sync"

type RWLocker struct {
	numReaders int        // number of readers holding lock
	hasWriter  bool       // if a writer requested lock
	c          *sync.Cond // conditional synchronization
}

func NewRWLocker() *RWLocker {
	return &RWLocker{
		numReaders: 0,
		hasWriter:  false,
		c:          sync.NewCond(&sync.Mutex{}),
	}
}

// Read lock
func (l *RWLocker) RLock() {
	l.c.L.Lock()
	// this check helps to prevent writer starvation
	for l.hasWriter {
		l.c.Wait()
	}
	l.numReaders += 1
	l.c.L.Unlock()
}

// Read unlock
func (l *RWLocker) RUnlock() {
	l.c.L.Lock()
	if l.numReaders >= 1 {
		l.numReaders -= 1
		// if this is the last reader, enable writer
		if l.numReaders == 0 {
			l.c.Broadcast()
		}
	}
	l.c.L.Unlock()
}

// Write lock
func (l *RWLocker) Lock() {
	l.c.L.Lock()

	// wait for writer
	for l.hasWriter {
		l.c.Wait()
	}
	l.hasWriter = true

	// wait for readers
	for l.numReaders != 0 {
		l.c.Wait()
	}

	l.c.L.Unlock()
}

// Write unlock
func (l *RWLocker) Unlock() {
	l.c.L.Lock()
	if l.hasWriter {
		l.hasWriter = false
		l.c.Broadcast()
	}
	l.c.L.Unlock()
}
