package segment

import (
	"context"
	"log"
)

// Destination interface has a blocking Process method, and Send method
type Destination interface {
	Process(ctx context.Context) error
	Send(ctx context.Context, message interface{}) error
	WithLogger(logger *log.Logger) Destination
}
