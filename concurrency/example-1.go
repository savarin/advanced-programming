/*
Problem: Print (mostly) the same number across all goroutines.

launched goroutine 3
launched goroutine 10
launched goroutine 10
launched goroutine 10
launched goroutine 10
launched goroutine 10
launched goroutine 10
launched goroutine 10
launched goroutine 10
launched goroutine 10

Solution: Introduce loop counter as argument to anonymous function.
https://github.com/savarin/advanced-programming/commit/0a242ca25171cf589ded4c699bba08198cb6ddc9#diff-d1b92de2d2079261b3abb4d275473f8ffefb9c723f42c8c56f72d385fb3abc7c
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("launched goroutine %d\n", i)
		}(i)
	}
	// Wait for goroutines to finish
	time.Sleep(time.Second)
}
