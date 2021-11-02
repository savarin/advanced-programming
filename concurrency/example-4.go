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
