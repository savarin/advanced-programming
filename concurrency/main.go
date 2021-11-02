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

type mutexIdService struct {
	counter uint64
	mu      sync.Mutex
}

func (s *mutexIdService) getNext() uint64 {
	s.mu.Lock()
	{
		s.counter += 1
	}
	s.mu.Unlock()

	return s.counter
}

type channelIdService struct {
	counter chan uint64
}

func newChannelIdService() *channelIdService {
	s := &channelIdService{
		counter: make(chan uint64),
	}

	go func() {
		var i uint64

		for {
			i += 1
			s.counter <- i
		}
	}()

	return s
}

func (s *channelIdService) getNext() uint64 {
	return <-s.counter
}

func run(wg *sync.WaitGroup, s idService, label string) {
	go func() {
		for i := 0; i < 999999; i++ {
			s.getNext()
		}

		fmt.Printf("1 %s:%d\n", label, s.getNext())
		wg.Done()
	}()

	go func() {
		for i := 0; i < 999999; i++ {
			s.getNext()
		}

		fmt.Printf("2 %s:%d\n", label, s.getNext())
		wg.Done()
	}()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(8)

	go func() {
		s1 := &basicIdService{}
		run(&wg, s1, "basic  ")
	}()
	go func() {
		s2 := &atomicIdService{}
		run(&wg, s2, "atomic ")
	}()
	go func() {
		s3 := &mutexIdService{}
		run(&wg, s3, "mutex  ")
	}()
	go func() {
		s4 := newChannelIdService()
		run(&wg, s4, "channel")
	}()

	fmt.Println("Waiting to finish...")
	wg.Wait()

	fmt.Println("Run complete, terminating.")
}
