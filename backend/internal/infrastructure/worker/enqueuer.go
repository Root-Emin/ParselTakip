// Package worker provides the Kafka-backed background job pool and an Enqueuer
// that publishes jobs onto the shared event bus.
package worker

import (
	"context"

	"github.com/masterfabric-go/masterfabric/internal/shared/events"
	"github.com/masterfabric-go/masterfabric/internal/shared/jobs"
)

// BusEnqueuer implements jobs.Enqueuer by publishing onto the event bus
// (Kafka in production, in-process for local/dev).
type BusEnqueuer struct {
	bus events.EventBus
}

var _ jobs.Enqueuer = (*BusEnqueuer)(nil)

// NewBusEnqueuer creates a new BusEnqueuer.
func NewBusEnqueuer(bus events.EventBus) *BusEnqueuer {
	return &BusEnqueuer{bus: bus}
}

// Enqueue publishes a job to the given topic for asynchronous processing.
func (e *BusEnqueuer) Enqueue(ctx context.Context, topic string, job interface{}) error {
	return e.bus.Publish(ctx, topic, job)
}
