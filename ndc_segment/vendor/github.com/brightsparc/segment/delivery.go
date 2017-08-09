package segment

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Create a summary to track delivery stream latency
	deliverySuccessCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "delivery_success_total",
		Help: "Delivery success total",
	}, []string{"stream"})
	deliveryFailureCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "delivery_failure_total",
		Help: "Delivery failure total",
	}, []string{"stream"})
	deliveryLatency = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "delivery_latency_seconds",
		Help:       "Delivery latency distributions",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"stream"})
	updatesSuccessCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "updates_success_total",
		Help: "Updates success total",
	}, []string{"stream"})
	updatesBlockedCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "updates_blocked_total",
		Help: "Updates blocked total",
	}, []string{"stream"})
	updatesLatency = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "updates_latency_seconds",
		Help:       "Updates latency distributions",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"stream"})
)

func init() {
	// Init prometheus metrics
	prometheus.MustRegister(deliverySuccessCounter)
	prometheus.MustRegister(deliveryFailureCounter)
	prometheus.MustRegister(deliveryLatency)
	prometheus.MustRegister(updatesSuccessCounter)
	prometheus.MustRegister(updatesBlockedCounter)
	prometheus.MustRegister(updatesLatency)
}

// DeliveryConfig contains configuration parameters including optional endpint
type DeliveryConfig struct {
	StreamEndpoint string        `json:"streamEndpoint,omitempty"`
	StreamRegion   string        `json:"streamRegion"`
	StreamName     string        `json:"streamName"`
	BatchSize      int           `json:"batchSize,omitempty"`
	FlushInterval  time.Duration `json:"flushInterval,omitempty"`
}

type Delivery struct {
	Logger        *log.Logger // Public logger that caller can override
	fh            *firehose.Firehose
	streamName    string
	size          int
	flushInterval time.Duration
	messages      chan interface{}
}

// NewDelivery creates a new delivery stream given configuration
func NewDelivery(config *DeliveryConfig) *Delivery {
	if config.StreamRegion == "" || config.StreamName == "" {
		log.Fatal("Require stream region and name")
	}
	if config.BatchSize <= 0 || config.BatchSize > 500 {
		config.BatchSize = 500
	}
	if config.FlushInterval == 0 {
		config.FlushInterval = time.Second * 30
	}

	// Block and initialize fh config on startup
	cfg := aws.NewConfig().WithRegion(config.StreamRegion)
	if config.StreamEndpoint != "" {
		cfg.WithEndpoint(config.StreamEndpoint)
	}
	sess := session.Must(session.NewSession(cfg))
	d := &Delivery{
		Logger:        log.New(os.Stderr, "", log.LstdFlags),
		fh:            firehose.New(sess, cfg),
		streamName:    config.StreamName,
		size:          config.BatchSize,
		flushInterval: config.FlushInterval,
	}

	return d
}

func (d *Delivery) WithLogger(logger *log.Logger) Destination {
	if logger != nil {
		d.Logger = logger
	}
	return d
}

// Process is a blocking call
func (d *Delivery) Process(ctx context.Context) error {
	log.Printf("Delivery connecting to %s...", d.fh.Endpoint)

	// Check the stream exists
	if stream, err := d.fh.DescribeDeliveryStream(&firehose.DescribeDeliveryStreamInput{
		DeliveryStreamName: aws.String(d.streamName),
	}); err != nil {
		return fmt.Errorf("Firehose stream error -- %v", err)
	} else {
		d.Logger.Println("Firehose stream", *stream.DeliveryStreamDescription.DeliveryStreamName)
	}

	// Create the async channel
	d.messages = make(chan interface{}, d.size*2)

	// Create the array to for batch of messages
	records := make([]*firehose.Record, d.size)

	send := func(i int) error {
		if i == 0 {
			d.Logger.Println("Nothing to send")
			return nil
		}

		t0 := time.Now()
		params := &firehose.PutRecordBatchInput{
			DeliveryStreamName: aws.String(d.streamName),
			Records:            records[:i],
		}
		resp, err := d.fh.PutRecordBatch(params)
		if err != nil {
			deliveryFailureCounter.WithLabelValues(d.streamName).Add(float64(i))
			d.Logger.Printf("Stream %s error sending %d: %s\n", d.streamName, i, err)
			return fmt.Errorf("Error sending to firehose -- %v", err)
		}

		// Log the succces, failed and latency metrics
		duration := time.Since(t0)
		deliveryFailureCounter.WithLabelValues(d.streamName).Add(float64(*resp.FailedPutCount))
		deliverySuccessCounter.WithLabelValues(d.streamName).Add(float64(i - int(*resp.FailedPutCount)))
		deliveryLatency.WithLabelValues(d.streamName).Observe(duration.Seconds())
		d.Logger.Printf("Stream %s sent %d (%d failed) in: %s\n", d.streamName, i, *resp.FailedPutCount, duration)
		return nil
	}

	d.Logger.Println("Starting processing")
	i := 0
	for {
		flush := false
		select {
		case message := <-d.messages:
			if data, err := json.Marshal(message); err != nil {
				return fmt.Errorf("Marshal error -- %v", err)
			} else {
				records[i] = &firehose.Record{
					Data: []byte(string(data) + "\n"), // Append newline after the json serialization
				}
				i++
			}
		case <-ctx.Done():
			// Sending remaining and return
			d.Logger.Println("Ending processing")
			return send(i)
		case <-time.After(d.flushInterval):
			if i > 0 {
				d.Logger.Printf("Flush after %s\n", d.flushInterval)
				flush = true
			}
		}
		if i == d.size || flush {
			// Send and reset index (records will be overwritten)
			send(i)
			i = 0
		}
	}
}

func (d *Delivery) Send(ctx context.Context, message interface{}) error {
	if d.messages == nil {
		return fmt.Errorf("Delivery destination not ready, check stream %q exists at %s", d.streamName, d.fh.Endpoint)
	}
	select {
	case d.messages <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
