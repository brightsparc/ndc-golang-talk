package main

import (
	"testing"
	"time"

	"github.com/brightsparc/fasttextgo"
)

func TestPredict(t *testing.T) {
	// Load model into memory
	t0 := time.Now()
	fasttextgo.LoadModel("model.bin")
	t.Logf("Model loaded in %s", time.Since(t0))

	// Predict
	s := "skills this talk is for anyone who wants to enjoy work more today happiness at work is no longer a luxury but an essential element for you to be productive engaged and happy i e to be the best person you can be  but for most of us work is work not a hobby nor a passion we suffer from smondays the moment when sunday stops feeling like a sunday as the anxiety of monday kicks in and this is not a sustainable way to live i was one of them in this talk i ll take you on a journey to help you understand what it is about work that is contributing to your unhappiness provide simple and effective tools that you can use to create a better work life for yourself and give you a massive dose of energy to kick you in the butt to make the change this will be a highly engaging and entertaining talk that is sure to challenge motivate and inspire"
	prob, label, err := fasttextgo.Predict(s)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Predict: %s (%f)\n", label, prob)
	}
}
