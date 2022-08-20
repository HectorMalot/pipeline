package pipeline

import "time"

// SendOrDrop passes an `Item any` from the `in <-chan Item` directly to the out `<-chan Item` unless the out channel is blocked.
// If the out channel is blocked, messages from `in <-chan Item` is sent to the `blocked` func instead.
func SendOrDrop[Item any](blocked func(Item), in <-chan Item) <-chan Item {
	out := make(chan Item)
	go func() {
		defer close(out)
		for i := range in {
			select {
			// When the channel isn't blocked, pass everything to the out chan
			case out <- i:

			// When the channel is blocked, pass the item to the blocked fun
			default:
				blocked(i)
			}
		}
	}()
	return out
}

// SendOrTimeout passes an `Item any` from the `in <-chan Item` directly to the out `<-chan Item` unless the out channel is blocked.
// If the out channel is blocked, it will wait for a duration of `timeout` and then pass the messages from `in <-chan Item`  to the `blocked` func instead.
func SendOrTimeout[Item any](timeout time.Duration, blocked func(Item), in <-chan Item) <-chan Item {
	out := make(chan Item)
	go func() {
		defer close(out)
		for i := range in {
			select {
			// When the channel isn't blocked, pass everything to the out chan
			case out <- i:

			// When the channel is blocked for 'timeout', pass the item to the blocked fun
			case <-time.After(timeout):
				blocked(i)
			}
		}
	}()
	return out
}
