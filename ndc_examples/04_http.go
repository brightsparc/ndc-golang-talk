package main // go run 05_http.go & curl http://localhost:8080/sydney

import (
	"fmt"
	"net/http"
)

func main() {
	// Inline func handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi NDC %s!", r.URL.Path[1:])
	})

	// Serve the endpoint
	http.ListenAndServe(":8080", nil)
}
