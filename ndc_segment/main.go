package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/brightsparc/segment"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ProjectConfig maps to json config file
type ProjectConfig struct {
	WriteKey        string                  `json:"writeKey"`
	ProjectId       string                  `json:"projectId"`
	Delivery        *segment.DeliveryConfig `json:"delivery"`
	ForwardEndpoint string                  `json:"forwardEndpoint"`
}

// GetProjectId returns the projectId from config file
func (c ProjectConfig) GetProjectId(writeKey string) string {
	if c.WriteKey == writeKey {
		return c.ProjectId
	}
	return ""
}

func loadConfig(filename string) ProjectConfig {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var config ProjectConfig
	if err := json.Unmarshal(raw, &config); err != nil {
		log.Fatalf("Error loading config %q -- %s", filename, err)
	}
	return config
}

var (
	version        = "1.0.0"
	apiHost        = flag.String("apiHost", ":3001", "Api Port")
	configFilename = flag.String("config", "./config.json", "Project Config")
	shutdown       = 5 * time.Second // The timeout for shutdown triggering (this needs to be long enough for events/updates to finish)
)

// Use the following command to run local firehose in docker:
// $ docker run -it -e SERVICES:firehose -p 4573:4573 -p 8080:8080 atlassianlabs/localstack

func main() {
	flag.Parse()

	// Declare stop signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Create the root and version endpoint
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "NDC API %s", version) // Pull version from tag
	})
	router.Handle("/metrics", promhttp.Handler()) // prometheus metrics endpoint
	vr := router.PathPrefix("/v1").Subrouter()    // Hard code version

	// Run the server and wait for complete
	h := &http.Server{Addr: *apiHost, Handler: router}
	go func() {
		log.Printf("NDC API %s listening on: %s\n", version, *apiHost)
		if err := h.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Load config, create a delivery destination, and segment forwarder destionation
	config := loadConfig(*configFilename)
	destinations := []segment.Destination{
		segment.NewDelivery(config.Delivery),
		segment.NewForwarder(config.ForwardEndpoint),
	}
	seg := segment.NewSegment(config.GetProjectId, destinations, vr)

	// Run waiting for cancel
	ctx, cancel := context.WithCancel(context.Background())
	seg.Run(ctx)

	// On stop signal, wait for quit to return
	<-stop
	cancel()

	log.Println("Shutting down web server...")
	ctx, cancel = context.WithTimeout(context.Background(), shutdown)
	defer cancel()

	h.Shutdown(ctx)
	log.Println("Server gracefully stopped")
}
