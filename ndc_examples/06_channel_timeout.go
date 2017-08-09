package main

import (
	"fmt"
	"time"
)

func main() {
	// Create a channel with empty struct
	ch := make(chan struct{})

	select {
	case _ = <-ch:
		fmt.Println("got result")
	case <-time.After(time.Second * 1):
		fmt.Println("timeout 1")
	}
}
