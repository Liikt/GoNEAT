package substrate

import (
	"sync"

	uuid "github.com/satori/go.uuid"
)

func genID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func merge(cs ...chan float64) chan float64 {
	var wg sync.WaitGroup
	out := make(chan float64)

	output := func(c chan float64) {
		out <- <-c
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
