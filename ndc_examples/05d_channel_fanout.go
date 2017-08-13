package main

import "fmt"

func square(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	for n := range square(2, 3) {
		fmt.Println(n) // 4, 9
	}
}
