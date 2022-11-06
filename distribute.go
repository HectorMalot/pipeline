package pipeline

// Distribute from 1 channel to many.
// Every subscribed channel receives every message
// Distribute blocks if the receiving channel
// is blocked. Use SendOrTimeout, SendOrDrop, or Buffer
// to prevent a blocked distributor
func Distribute[Item any](in <-chan Item, size int) []<-chan Item {
	out := make([]<-chan Item, size)
	chans := make([]chan Item, size)
	for i := 0; i < size; i++ {
		ch := make(chan Item, 0)
		out[i] = ch
		chans[i] = ch
	}

	go func() {
		for v := range in {
			for _, ch := range chans {
				ch <- v
			}
		}
		for _, ch := range chans {
			close(ch)
		}
	}()

	return out
}
