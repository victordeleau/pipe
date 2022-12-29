package lib

import (
	"github.com/victordeleau/pipe/pkg"
	"sync"
)

func Fanin[T any](inputs ...pipe.Chan[T]) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)
	fan := func(c pipe.Chan[T]) {
		for n := range c.Buffer {
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
