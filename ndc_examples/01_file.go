package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.Open("talk.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
}
