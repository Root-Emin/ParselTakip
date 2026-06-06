package worker

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	domainStorage "github.com/masterfabric-go/masterfabric/internal/domain/storage"
	"github.com/masterfabric-go/masterfabric/internal/shared/events"
	"github.com/masterfabric-go/masterfabric/internal/shared/jobs"
)

// Pool registers background job handlers on the event bus with bounded retry and
// dead-letter semantics. It works with both the Kafka and in-process buses.
type Pool struct {
	bus        events.EventBus
	logger     *slog.Logger
	maxRetries int
}

// NewPool creates a worker pool. maxRetries < 1 defaults to 3.
func NewPool(bus events.EventBus, logger *slog.Logger, maxRetries int) *Pool {
	if maxRetries < 1 {
		maxRetries = 3
	}
	return &Pool{bus: bus, logger: logger, maxRetries: maxRetries}
}

// Handler processes a single job event.
type Handler func(ctx context.Context, e events.Event) error

// Register subscribes a handler to a job topic. The handler is retried with a
// small linear backoff; once retries are exhausted the job is routed to the
// dead-letter topic and acknowledged so the queue can make progress.
func (p *Pool) Register(topic string, handler Handler) {
	p.bus.Subscribe(topic, func(ctx context.Context, e events.Event) error {
		var lastErr error
		for attempt := 1; attempt <= p.maxRetries; attempt++ {
			if lastErr = handler(ctx, e); lastErr == nil {
				return nil
			}
			p.logger.Warn("job attempt failed", "topic", topic, "attempt", attempt, "max", p.maxRetries, "error", lastErr)
			if attempt < p.maxRetries {
				time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
			}
		}
		p.deadLetter(ctx, topic, e, lastErr)
		return nil
	})
}

func (p *Pool) deadLetter(ctx context.Context, topic string, e events.Event, cause error) {
	payload, _ := json.Marshal(rawEvent(e))
	msg := "unknown error"
	if cause != nil {
		msg = cause.Error()
	}
	_ = p.bus.Publish(ctx, jobs.TopicDeadLetter, jobs.DeadLetter{
		OriginalTopic: topic,
		Payload:       string(payload),
		Error:         msg,
		FailedAt:      time.Now().UTC().Format(time.RFC3339),
	})
	p.logger.Error("job dead-lettered", "topic", topic, "error", msg)
}

// RegisterDefaults wires the standard background job handlers.
func (p *Pool) RegisterDefaults(storage domainStorage.ObjectStorage) {
	// Post-upload document processing: verify the object landed in storage
	// (integrity check) before marking it done. Heavy work (AV scan, OCR) hooks here.
	p.Register(jobs.TopicDocumentProcess, func(ctx context.Context, e events.Event) error {
		job, err := decode[jobs.DocumentProcessJob](e)
		if err != nil {
			return err
		}
		if storage != nil && job.StorageKey != "" {
			if _, err := storage.Stat(ctx, job.StorageKey); err != nil {
				return err
			}
		}
		p.logger.Info("document processed", "document_id", job.DocumentID, "checksum", job.Checksum)
		return nil
	})

	// Notification dispatch (push/email/SMS delivery happens off the request path).
	p.Register(jobs.TopicNotificationDispatch, func(ctx context.Context, e events.Event) error {
		job, err := decode[jobs.NotificationDispatchJob](e)
		if err != nil {
			return err
		}
		p.logger.Info("notification dispatched", "notification_id", job.NotificationID, "channel", job.Channel)
		return nil
	})

	// Outbound email.
	p.Register(jobs.TopicEmailSend, func(ctx context.Context, e events.Event) error {
		job, err := decode[jobs.EmailSendJob](e)
		if err != nil {
			return err
		}
		p.logger.Info("email sent", "to", job.To, "subject", job.Subject)
		return nil
	})
}

// rawEvent returns a JSON-able payload regardless of bus transport.
func rawEvent(e events.Event) interface{} {
	if env, ok := e.(*events.Envelope); ok {
		return json.RawMessage(env.Data)
	}
	return e
}

// decode extracts a typed job from an event, transparently handling both
// in-process delivery (the original Go struct) and Kafka delivery (*Envelope).
func decode[T any](e events.Event) (T, error) {
	var out T
	switch v := e.(type) {
	case *events.Envelope:
		if err := json.Unmarshal(v.Data, &out); err != nil {
			return out, err
		}
	default:
		b, err := json.Marshal(e)
		if err != nil {
			return out, err
		}
		if err := json.Unmarshal(b, &out); err != nil {
			return out, err
		}
	}
	return out, nil
}
