package segment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Forwarder struct {
	Logger   *log.Logger // Public logger that caller can override
	endpoint string
}

// NewDelivery creates a new delivery stream given configuration
func NewForwarder(endpoint string) *Forwarder {
	return &Forwarder{
		Logger:   log.New(os.Stderr, "", log.LstdFlags),
		endpoint: endpoint,
	}
}

func (f *Forwarder) WithLogger(logger *log.Logger) Destination {
	if logger != nil {
		f.Logger = logger
	}
	return f
}

func (f *Forwarder) Process(ctx context.Context) error {
	log.Println("Started forwarder")
	<-ctx.Done() // Block on context, since work is done in Send
	return nil
}

func (f *Forwarder) Send(ctx context.Context, message interface{}) error {
	// Cast this message to a segment event
	m, ok := message.(SegmentEvent)
	if !ok {
		return fmt.Errorf("Expected Segment Event")
	}
	batch := SegmentBatch{
		MessageId: m.MessageId,
		Timestamp: m.Timestamp,
		SentAt:    m.SentAt,
		Context:   m.Context,
		Messages:  []SegmentMessage{m.SegmentMessage},
	}
	b, err := json.Marshal(batch)
	if err != nil {
		return err
	}

	// Create the request for the specific type
	req, err := http.NewRequest("POST", f.endpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}
	req.Header.Add("User-Agent", "analytics-go (version: 2.1.0)")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(b)))
	req.SetBasicAuth(m.WriteKey, "")

	// Send request
	err = httpDo(ctx, req, func(res *http.Response, err error) error {
		if err != nil {
			return fmt.Errorf("error sending request: %s", err)
		}
		defer res.Body.Close()
		if res.StatusCode < 400 {
			return nil
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %s", err)
		}
		return fmt.Errorf("response %s: %d – %s", res.Status, res.StatusCode, string(body))
	})

	if err == nil {
		f.Logger.Println("Forwarded", m.MessageId)
	}
	return err
}

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}
