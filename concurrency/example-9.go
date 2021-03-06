/*
Problem: Deadlock due to one goroutine acquires StateManager lock first and Consumer lock second,
and another goroutine in the opposite ordering.

Solution: Change lock acquisition ordering so that both acquire the Statemanager lock first and
Consumer lock second.
https://github.com/savarin/advanced-programming/commit/4a234a8aa3f2e020292a33cf55d20232543045df#diff-d1b92de2d2079261b3abb4d275473f8ffefb9c723f42c8c56f72d385fb3abc7c
*/
package main

import (
	"fmt"
	"sync"
)

type Consumer struct {
	id   int
	s    *StateManager
	lock *sync.RWMutex
}

func NewConsumer(id int, s *StateManager) *Consumer {
	return &Consumer{
		id:   id,
		s:    s,
		lock: &sync.RWMutex{},
	}
}

func (c *Consumer) GetState() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return fmt.Sprintf("<GetState result for consumer %d>", c.id)
}

func (c *Consumer) Terminate() {
	c.lock.Lock()
	{
		// You can imagine that this internal cleanup mutates the state
		// of the Consumer
		fmt.Printf("Performing internal cleanup for consumer %d\n", c.id)
	}
	c.lock.Unlock()
}

type StateManager struct {
	lock      *sync.RWMutex
	consumers map[int]*Consumer
}

func NewStateManager(numConsumers int) *StateManager {
	s := &StateManager{
		lock:      &sync.RWMutex{},
		consumers: make(map[int]*Consumer),
	}
	for i := 0; i < numConsumers; i++ {
		s.consumers[i] = NewConsumer(i, s)
	}
	return s
}

func (s *StateManager) AddConsumer(c *Consumer) {
	s.lock.Lock()
	{
		s.consumers[c.id] = c
	}
	s.lock.Unlock()
}

func (s *StateManager) RemoveConsumer(c *Consumer) {
	s.lock.Lock()
	{
		delete(s.consumers, c.id)
		c.Terminate()
	}
	s.lock.Unlock()
}

func (s *StateManager) GetConsumer(id int) *Consumer {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.consumers[id]
}

func (s *StateManager) PrintState() {
	fmt.Println("Started PrintState")

	s.lock.RLock()
	{
		for _, consumer := range s.consumers {
			fmt.Println(consumer.GetState())

		}
	}
	s.lock.RUnlock()

	fmt.Println("Done with PrintState")
}

func main() {
	s := NewStateManager(10)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		c := s.GetConsumer(0)
		s.RemoveConsumer(c)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		s.PrintState()
	}()

	wg.Wait()
}
