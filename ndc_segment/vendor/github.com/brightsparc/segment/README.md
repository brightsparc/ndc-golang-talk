# Segment

[Segment.io](https://segment.com/) compatible server written in go.

## Introduction

Segment is a cloud based analytics platform for tracking events from your application.  

It has a well designed [specific](https://segment.com/docs/spec/) that supports APIs:
* [Page](https://segment.com/docs/spec/page/): what web page are they on?
* [Track](https://segment.com/docs/spec/track/): what are they doing?
* [Identify](https://segment.com/docs/spec/identify/): who is the customer?
* [Alias](https://segment.com/docs/spec/alias/): what was their past identity?
* [Group](https://segment.com/docs/spec/group/): what account or organization are they part of?
* [Screen](https://segment.com/docs/spec/screen/): what app screen are they on?

This go library implements endpoints for all of these APIs.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

This library uses the [gorilla mux](https://github.com/gorilla/mux) library for attaching http handlers, and prometheus for monitoring.

```
go get -u github.com/gorilla/mux
go get -u github.com/prometheus/client_golang
```

### Installing

This library is installed as package `segment`.

```
go get -u  github.com/brightsparc/segment
```

## Examples

Create a new Segment listener by providing a function to return projectId from writeKey.  For unknown writeKey values, return empty string to have endpoint return 400 back request. Configure one or more destinations, this example includes forwarded to segment cloud, and firehose stream.

```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/brightsparc/segment"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	projectId := func(writeKey string) string {
		return "xxxx" // TODO: Match this with your writeKey for authorisation
	}
	destinations := []segment.Destination{
		segment.NewForwarder("https://api.segment.io/v1/batch"),
		segment.NewDelivery(&segment.DeliveryConfig{
			StreamRegion: "us-west-2",
			StreamName:   "stream-name", // Must exist
		}),
	}

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler()) // prometheus metrics endpoint
	sr := router.PathPrefix("/v1").Subrouter()    // Will create endpoints /v1/batch etc
	seg := segment.NewSegment(projectId, destinations, sr)
	go seg.Run(context.Background())

	log.Println("Listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
```

## Implementation Details

### Send messages

The segment `Send` method will execute `Send` method on each destination in order, and return on error.  It is recommended to implement a queue as per the `Delivery` process, the `Forwarder` should only be used for testing.

### Background process

* The segment `Run` method processes destinations on seperate go routines, blocking until the context is done.
* The firehose `Delivery` process batches up to 500 messages, sending them at every 30 seconds by default.

### Logging

The `Segment` class will log to standard error by default, but can be configured by the `Logger` property.

### Monitoring

The [prometheus](https://github.com/prometheus/client_golang) client is enabled to return http and delivery metrics.  

## Authors

* **Julian Bright** - [brightsparc](https://github.com/brightsparc/)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
