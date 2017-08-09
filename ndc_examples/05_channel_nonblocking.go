package main

import (
	"fmt"
)

func main() {
	ch := make(chan string, 1)
	select {
	case ch <- "buffered":
		fmt.Println("message 1 sent")
	}
	select {
	case ch <- "non-blocking":
		fmt.Println("message 2 sent")
	default:
		fmt.Println("no message sent")
	}
}
