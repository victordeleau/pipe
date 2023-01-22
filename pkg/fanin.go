package pipe

import "sync"

func Fanin[T any](inputs ...Channel[T]) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)
	fan := func(c Channel[T]) {
		for n := range *c.Buffer {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(inputs))
	for _, c := range inputs {
		go fan(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
