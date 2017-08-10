package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/brightsparc/fasttextgo"
)

// Speaker request
type Speaker struct {
	Name     string `json:"name"`
	Tagline  string `json:"tagline"`
	Image    string `json:"image"`
	Preamble string `json:"preamble"`
	Talk     struct {
		Title    string   `json:"title"`
		Tags     []string `json:"tags"`
		Preamble string   `json:"preamble"`
		Body     string   `json:"body"`
	} `json:"talk"`
}

// Prediction Response
type Prediction struct {
	Prob  float32 `json:"prob"`
	Label string  `json:"label"`
}

func predict(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()

	var speaker Speaker
	err := json.NewDecoder(r.Body).Decode(&speaker)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get prob and label, return results
	var pred Prediction
	pred.Prob, pred.Label, err = fasttextgo.Predict(speaker.Talk.Title + " " + speaker.Talk.Body)
	if err != nil {
		http.Error(w, "Predict error", http.StatusInternalServerError)
	} else {
		log.Printf("Predict in %s\n", time.Since(t0))
		json.NewEncoder(w).Encode(pred)
	}
}

var (
	version       = "1.0.0"
	apiHost       = flag.String("apiHost", ":3002", "Api Port")
	modelFilename = flag.String("model", "../ndc_scraper/model.bin", "Model path")
)

func main() {
	flag.Parse()

	// Load model into memory
	t0 := time.Now()
	log.Printf("Loading model %s\n", *modelFilename)
	fasttextgo.LoadModel(*modelFilename)
	log.Printf("Model loaded in %s\n", time.Since(t0))

	// Add handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "NDC Predict API %s", version)
	})
	http.HandleFunc("/predict", predict)

	// Serve API
	log.Printf("NDC Predict API %s listening on: %s\n", version, *apiHost)
	log.Fatal(http.ListenAndServe(*apiHost, nil))
}
