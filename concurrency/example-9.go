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
	defer c.lock.Unlock()

	// You can imagine that this internal cleanup mutates the state
	// of the Consumer
	fmt.Printf("Performing internal cleanup for consumer %d\n", c.id)

	c.s.RemoveConsumer(c.id)
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
	defer s.lock.Unlock()

	s.consumers[c.id] = c
}

func (s *StateManager) RemoveConsumer(id int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.consumers, id)
}

func (s *StateManager) GetConsumer(id int) *Consumer {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.consumers[id]
}

func (s *StateManager) PrintState() {
	s.lock.RLock()
	defer s.lock.RUnlock()

	fmt.Println("Started PrintState")
	for _, consumer := range s.consumers {
		fmt.Println(consumer.GetState())
	}
	fmt.Println("Done with PrintState")
}

func main() {
	s := NewStateManager(10)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		c := s.GetConsumer(0)
		c.Terminate()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		s.PrintState()
	}()

	wg.Wait()
}
