package main // cat talk.json | go run 02_json.go

import (
	"encoding/json"
	"fmt"
	"os"
)

type Speaker struct {
	Talk struct {
		Title string `json:"title"`
	} `json:"talk"`
}

func main() {
	decoder := json.NewDecoder(os.Stdin)
	var speaker Speaker
	if err := decoder.Decode(&speaker); err == nil {
		fmt.Println(speaker.Talk.Title)
	}
}
