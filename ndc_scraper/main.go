package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/brightsparc/fasttextgo"
	flags "github.com/jessevdk/go-flags"
)

/*
# Train model
time ../../fasttextgo/fasttext supervised -input train.txt -output model -wordNgrams 2 -lr 0.1 -epoch 1000

# Test model between c and go
time ../../fasttextgo/fasttext predict-prob model.bin test.txt 1 | head -n 1
time go run main.go predict-prob model.bin test.txt 1 | head -n 1
*/

func main() {
	var opts = struct {
		Verbose    bool `short:"v"`
		Positional struct {
			Command   string
			ModelFile string
			TestFile  string
			Count     int
		} `positional-args:"yes" required:"yes"`
	}{}

	var parser = flags.NewParser(&opts, flags.IgnoreUnknown)
	if _, err := parser.Parse(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Args", opts)
	}

	// Load model
	t0 := time.Now()
	fasttextgo.LoadModel(opts.Positional.ModelFile)
	log.Printf("Model loaded in %s\n", time.Since(t0))

	switch opts.Positional.Command {
	case "predict", "predict-prob":
		// Load test
		f, err := os.Open(opts.Positional.TestFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		r := csv.NewReader(f)
		r.Comma = '\t' // Split on any string?
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			prob, label, err := fasttextgo.Predict(record[1])
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Print(label)
				if opts.Positional.Command == "predict-prob" {
					fmt.Printf(" %f\n", prob)
				} else {
					fmt.Println()
				}
			}
		}
	default:
		log.Fatalf("Command %s not supported", opts.Positional.Command)
	}
}
