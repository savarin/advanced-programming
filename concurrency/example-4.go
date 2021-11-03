/*
Problem: No initialization takes place because the receive, which is a blocking operation, is in the
goroutine and the send outside.

Solution: Swap send and receive operation.
https://github.com/savarin/advanced-programming/commit/e83f10a88a75659ee5638228e679d115e6bea92e#diff-d1b92de2d2079261b3abb4d275473f8ffefb9c723f42c8c56f72d385fb3abc7c
*/
package main

import (
	"fmt"
)

func main() {
	done := make(chan struct{}, 1)
	go func() {
		fmt.Println("performing initialization...")
		done <- struct{}{}
	}()

	<-done
	fmt.Println("initialization done, continuing with rest of program")
}
