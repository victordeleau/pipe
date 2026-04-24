package pipe

import "sync"

// Fanin merges multiple receive-only streams into one. Each input must
// eventually be closed (or the goroutine for that input will not finish).
// For pipe ReceiveChannel values, pass ch.ReadChan() per input.
func Fanin[T any](inputs ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)
	fan := func(ch <-chan T) {
		for n := range ch {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(inputs))
	for _, ch := range inputs {
		go fan(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
