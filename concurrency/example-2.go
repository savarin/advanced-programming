/*
Problem: Deadlock due to channel having nil value.

Solution: Declare channel with make.
https://github.com/savarin/advanced-programming/commit/162f366501d4de3c687a2a37f4e94c6cb7e7db02#diff-d1b92de2d2079261b3abb4d275473f8ffefb9c723f42c8c56f72d385fb3abc7c
*/
package main

import (
	"fmt"
)

const numTasks = 3

func main() {
	done := make(chan struct{})
	for i := 0; i < numTasks; i++ {
		go func() {
			fmt.Println("running task...")

			// Signal that task is done
			done <- struct{}{}
		}()
	}

	// Wait for tasks to complete
	for i := 0; i < numTasks; i++ {
		<-done
	}
	fmt.Printf("all %d tasks done!\n", numTasks)
}
