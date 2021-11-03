/*
Problem: Since there are 3 goroutines sending to a non-buffered channel, only one of the goroutines
can perform the send operation. This causes the number of goroutines run to increase over time as
the send-receive connection is not closed.

Solution: Change to a channel with a buffer of size 3.
https://github.com/savarin/advanced-programming/commit/6c5faa59d180cca5f05012105d1d120106cd9234#diff-d1b92de2d2079261b3abb4d275473f8ffefb9c723f42c8c56f72d385fb3abc7c
*/
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

var responses = []string{
	"200 OK",
	"402 Payment Required",
	"418 I'm a teapot",
}

func randomDelay(maxMillis int) time.Duration {
	return time.Duration(rand.Intn(maxMillis)) * time.Millisecond
}

func query(endpoint string) string {
	// Simulate querying the given endpoint
	delay := randomDelay(100)
	time.Sleep(delay)

	i := rand.Intn(len(responses))
	return responses[i]
}

// Query each of the mirrors in parallel and return the first
// response (this approach increases the amount of traffic but
// significantly improves "tail latency")
func parallelQuery(endpoints []string) string {
	results := make(chan string, len(endpoints))
	for i := range endpoints {
		go func(i int) {
			results <- query(endpoints[i])
		}(i)
	}
	return <-results
}

func main() {
	var endpoints = []string{
		"https://fakeurl.com/endpoint",
		"https://mirror1.com/endpoint",
		"https://mirror2.com/endpoint",
	}

	// Simulate long-running server process that makes continuous queries
	for {
		fmt.Println(parallelQuery(endpoints))
		delay := randomDelay(100)
		time.Sleep(delay)
		fmt.Println("number of goroutines: ", runtime.NumGoroutine())
	}
}
