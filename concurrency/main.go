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

		fmt.Printf("%s  1:%d\n", label, s.getNext())
		wg.Done()
	}()

	go func() {
		for i := 0; i < 99999; i++ {
			s.getNext()
		}

		fmt.Printf("%s  2:%d\n", label, s.getNext())
		wg.Done()
	}()
}

func runWithMutex(wg *sync.WaitGroup, mu *sync.Mutex, s idService) {
	var c uint64

	go func() {
		for i := 0; i < 99999; i++ {
			mu.Lock()
			{
				s.getNext()
			}
			mu.Unlock()
		}

		mu.Lock()
		{
			c = s.getNext()
		}
		mu.Unlock()

		fmt.Printf("mutex   1:%d\n", c)
		wg.Done()
	}()

	go func() {
		for i := 0; i < 99999; i++ {
			mu.Lock()
			{
				s.getNext()
			}
			mu.Unlock()
		}

		mu.Lock()
		{
			c = s.getNext()
		}
		mu.Unlock()

		fmt.Printf("mutex   2:%d\n", c)
		wg.Done()
	}()
}

func runWithChannels(wg *sync.WaitGroup, s idService) {
	ch1 := make(chan uint64)
	ch2 := make(chan uint64)

	go func() {
		for i := 0; i < 100000; i++ {
			ch1 <- s.getNext()
		}

		wg.Done()
	}()

	go func() {
		for i := 0; i < 100000; i++ {
			ch2 <- s.getNext()
		}

		wg.Done()
	}()

	var c uint64

	for i := 0; i < 100000; i++ {
		c = <-ch1
	}

	fmt.Printf("channel 1:%d\n", c)

	for i := 0; i < 100000; i++ {
		c = <-ch2
	}

	fmt.Printf("channel 2:%d\n", c)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(8)

	var mu sync.Mutex

	s1 := &basicIdService{}
	s2 := &atomicIdService{}
	s3 := &basicIdService{}
	s4 := &basicIdService{}

	run(&wg, s1, "basic ")
	run(&wg, s2, "atomic")
	runWithMutex(&wg, &mu, s3)
	runWithChannels(&wg, s4)

	fmt.Println("Waiting to finish...")
	wg.Wait()

	fmt.Println("Run complete, terminating.")
}
