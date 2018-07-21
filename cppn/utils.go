package cppn

import (
	"sync"

	"github.com/satori/go.uuid"
)

var (
	innovation      = 0
	innovationMutex = new(sync.Mutex)
)

func keysNeuron(m map[string]*neuron) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return
}

func keyGene(m map[string]*gene) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return
}

func indexGene(m map[int]*gene) (keys []int) {
	for k := range m {
		keys = append(keys, k)
	}
	return
}

func newInnovation() int {
	innovationMutex.Lock()
	defer innovationMutex.Unlock()
	innovation++
	return innovation
}

func containsInt(arr []int, i int) bool {
	for _, x := range arr {
		if i == x {
			return true
		}
	}
	return false
}

func containsString(arr []string, i string) bool {
	for _, x := range arr {
		if i == x {
			return true
		}
	}
	return false
}

func findMissing(one, two []int, c chan int) {
	counter := 0
	for _, g := range one {
		if !containsInt(two, g) {
			counter++
		}
	}
	c <- counter
}

func matchingInts(one, two []int) []int {
	ret := make([]int, 0)

	for _, g := range one {
		if containsInt(two, g) {
			ret = append(ret, g)
		}
	}

	return ret
}

func truncateValue(value float64) float64 {
	if value < -3.0 {
		return -3.0
	} else if value > 3.0 {
		return 3.0
	}

	return value
}

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
