package main

import (
	"testing"
)

func BenchmarkBasic(b *testing.B) {
	s := new(basicIdService)
	for i := 0; i < b.N; i++ {
		_ = s.getNext()
	}
}

func BenchmarkAtomic(b *testing.B) {
	s := new(atomicIdService)
	for i := 0; i < b.N; i++ {
		_ = s.getNext()
	}
}

func BenchmarkMutex(b *testing.B) {
	s := new(mutexIdService)
	for i := 0; i < b.N; i++ {
		_ = s.getNext()
	}
}

func BenchmarkChannel(b *testing.B) {
	s := newChannelIdService()
	for i := 0; i < b.N; i++ {
		_ = s.getNext()
	}
}
