package main

func main() {
	// create new channel of type int
	ch := make(chan int, 1)

	// start new anonymous goroutine
	go func() {
		// send 42 to channel
		ch <- 42
	}()

	// read from channel
	<-ch
}
