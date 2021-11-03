/*
Problem: Deadlock due attempt to acquire read lock already held.

Solution: Remove need to acquire lock when lock already held.
https://github.com/savarin/advanced-programming/commit/93bfc23b2639f428134a1ecc032cba10e70e1c47#diff-d1b92de2d2079261b3abb4d275473f8ffefb9c723f42c8c56f72d385fb3abc7c
*/
package main

import (
	"fmt"
	"sync"
	"time"
)

type dbService struct {
	lock       *sync.RWMutex
	connection string
}

func newDbService(connection string) *dbService {
	return &dbService{
		lock:       &sync.RWMutex{},
		connection: connection,
	}
}

func (d *dbService) logState() {
	d.lock.RLock()
	{
		d.printState()
	}
	d.lock.RUnlock()
}

func (d *dbService) printState() {
	fmt.Printf("connection %q is healthy\n", d.connection)
}

func (d *dbService) takeSnapshot() {
	d.lock.RLock()
	{
		fmt.Printf("Taking snapshot over connection %q\n", d.connection)

		// Simulate slow operation
		time.Sleep(time.Second)

		d.printState()
	}
	d.lock.RUnlock()
}

func (d *dbService) updateConnection(connection string) {
	d.lock.Lock()
	{
		d.connection = connection
	}
	d.lock.Unlock()
}

func main() {
	d := newDbService("127.0.0.1:3001")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		d.takeSnapshot()
	}()

	// Simulate other DB accesses
	time.Sleep(200 * time.Millisecond)

	wg.Add(1)
	go func() {
		defer wg.Done()

		d.updateConnection("127.0.0.1:8080")
	}()

	wg.Wait()
}
