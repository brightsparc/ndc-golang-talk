package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brightsparc/fasttextgo"
	flags "github.com/jessevdk/go-flags"
)

/*
# Train model
time ../../fasttextgo/fasttext supervised -input train.txt -output model -wordNgrams 2 -lr 0.1 -epoch 1000

# Test models
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
	_, err := parser.Parse()
	if err != nil {
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
		f, err := os.Open(opts.Positional.TestFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			prob, label, err := fasttextgo.Predict(scanner.Text())
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
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading input: %s\n", err)
		}
	default:
		log.Fatalf("Command %s not supported", opts.Positional.Command)
	}
}
