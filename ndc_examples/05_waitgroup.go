package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	do := func(i int) {
		fmt.Println(i * i) // Do work
		wg.Done()
	}
	wg.Add(2)
	go do(2)
	go do(3)

	fmt.Println("Wating...")
	wg.Wait()
}
