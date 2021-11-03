/*
Problem: (1) WaitGroup does not wait for all threads to complete, and (2) lock passed by value.

Solution: (1) Increase WaitGroup to number of goroutines, and (2) pass lock by reference.
https://github.com/savarin/advanced-programming/commit/549993f69b5a88933a5e3e046ff42f3488732011#diff-d1b92de2d2079261b3abb4d275473f8ffefb9c723f42c8c56f72d385fb3abc7c
*/
package main

import (
	"fmt"
	"sync"
)

const (
	numGoroutines = 100
	numIncrements = 100
)

type counter struct {
	count int
}

func safeIncrement(lock *sync.Mutex, c *counter) {
	lock.Lock()
	{
		c.count += 1
	}
	lock.Unlock()
}

func main() {
	var globalLock sync.Mutex
	c := &counter{
		count: 0,
	}

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < numIncrements; j++ {
				safeIncrement(&globalLock, c)
			}
		}()
	}

	wg.Wait()
	fmt.Println(c.count)
}
