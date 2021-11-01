package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type idService interface {
	// Returns values in ascending order; it should be safe to call getNext() concurrently without
	// any additional synchronization.
	getNext() uint64
}

type basicIdService struct {
	counter uint64
}

func (s *basicIdService) getNext() uint64 {
	s.counter += 1
	return s.counter
}

type atomicIdService struct {
	counter uint64
}

func (s *atomicIdService) getNext() uint64 {
	atomic.AddUint64(&s.counter, 1)
	return s.counter
}

func run(wg *sync.WaitGroup, s idService, label string) {
	go func() {
		for i := 0; i < 99999; i++ {
			s.getNext()
		}

		fmt.Printf("%s 1:%d\n", label, s.getNext())
		wg.Done()
	}()

	go func() {
		for i := 0; i < 99999; i++ {
			s.getNext()
		}

		fmt.Printf("%s 2:%d\n", label, s.getNext())
		wg.Done()
	}()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(4)

	s1 := &basicIdService{}
	s2 := &atomicIdService{}

	run(&wg, s1, "basic")
	run(&wg, s2, "atomic")

	fmt.Println("Waiting to finish...")
	wg.Wait()

	fmt.Println("Run complete, terminating.")
}
