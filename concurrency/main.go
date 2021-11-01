package main

import (
	"fmt"
	"sync"
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

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	s := &basicIdService{}

	go func() {
		for i := 0; i < 1000; i++ {
			s.getNext()
		}

		fmt.Println("goroutine 1:", s.counter)
		wg.Done()
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			s.getNext()
		}

		fmt.Println("goroutine 2:", s.counter)

		wg.Done()
	}()

	fmt.Println("Waiting to finish...")
	wg.Wait()

	fmt.Println("Run complete, terminating.")
}
