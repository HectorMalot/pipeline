package pipeline

import (
	"sync"
)

// Distributor forwards received messages to all receivers
type Distributor[T any] struct {
	in    <-chan T
	out   []chan T
	mutex sync.RWMutex
}

// Distribute from 1 channel to many.
// Every subscribed channel receives every message
// Distribute blocks if the receiving channel
// is blocked. Use SendOrTimeout, SendOrDrop, or Buffer
// to prevent a blocked distributor
func NewDistributor[T any](in <-chan T) *Distributor[T] {
	d := &Distributor[T]{in: in, out: make([]chan T, 0)}
	go func() {
		for v := range in {
			d.mutex.RLock()
			for _, ch := range d.out {
				ch <- v
			}
			d.mutex.RUnlock()
		}
		d.closeAllChannels()
	}()
	return d
}

// Get a new output channel of the distributor
func (d *Distributor[T]) NewChan() <-chan T {
	ch := make(chan T)
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.out = append(d.out, ch)
	return ch
}

// Remove an output channel of the distributor
func (d *Distributor[T]) RemoveChan(channel chan T) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	for i, ch := range d.out {
		if ch == channel {
			d.out[len(d.out)-1], d.out[i] = nil, d.out[len(d.out)-1]
			d.out = d.out[:len(d.out)-1]
			close(channel)
			return true
		}
	}
	return false
}

func (d *Distributor[T]) closeAllChannels() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	for _, ch := range d.out {
		close(ch)
	}
	d.out = make([]chan T, 0)
}
