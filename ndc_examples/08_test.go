package main // go test -v -bench=. 06_test.go

import "testing"

func TestSucceed(t *testing.T) {
	t.Log("Do something useful")
}

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

func BenchmarkFib10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fib(10) // run the Fib function b.N times
	}
}
