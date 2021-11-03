/*
Problem: Deadlock due to read lock acquisition attempted while already holding the lock.

Solution: Introduce logging without locks, for use when lock already held.
https://github.com/savarin/advanced-programming/commit/ce9071f0fd7bb207df07282e1ef5ef4fcae700e6#diff-d1b92de2d2079261b3abb4d275473f8ffefb9c723f42c8c56f72d385fb3abc7c
*/
package main

import (
	"fmt"
	"sync"
)

type coordinator struct {
	lock   sync.RWMutex
	leader string
}

func newCoordinator(leader string) *coordinator {
	return &coordinator{
		lock:   sync.RWMutex{},
		leader: leader,
	}
}

func (c *coordinator) logState() {
	c.lock.RLock()
	{
		c.printState()
	}
	c.lock.RUnlock()
}

// Assumes lock already acquired.
func (c *coordinator) printState() {
	fmt.Printf("leader = %q\n", c.leader)
}

func (c *coordinator) setLeader(leader string, shouldLog bool) {
	c.lock.Lock()
	{
		c.leader = leader

		if shouldLog {
			c.printState()
		}
	}
	c.lock.Unlock()
}

func main() {
	c := newCoordinator("us-east")
	c.logState()
	c.setLeader("us-west", true)
}
