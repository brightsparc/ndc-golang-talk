package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brightsparc/fasttextgo"
)

// Session request
type Session struct {
	Speaker struct {
		Name     string `json:"name"`
		Tagline  string `json:"tagline"`
		Image    string `json:"image"`
		Preamble string `json:"preamble"`
	} `json:"speaker"`
	Talk struct {
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

const prefix = "__label__"

func predict(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	defer log.Printf("Predict in %s\n", time.Since(t0))

	// Decode the session request
	var session Session
	err := json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get top K parameter from query string
	var k int
	top := r.URL.Query().Get("top")
	if k, err = strconv.Atoi(top); err != nil {
		k = 3
	}

	// Get prob and label, return results
	probs, lables, err := fasttextgo.PredictK(session.Talk.Title+" "+session.Talk.Body, k)
	if err != nil {
		http.Error(w, "Predict error", http.StatusInternalServerError)
		return
	}

	// Return a list of predictions
	var preds []Prediction
	for i := range probs {
		preds = append(preds, Prediction{
			Prob:  probs[i],
			Label: strings.TrimPrefix(lables[i], prefix),
		})
	}
	json.NewEncoder(w).Encode(preds)
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
