// Package jobs defines the background job contracts (topics, payloads and the
// Enqueuer interface) used to offload heavy or high-volume work from the request
// path onto the Kafka-backed worker pool.
package jobs

import "context"

// Background job topics. These are dedicated Kafka topics consumed by the worker
// pool, separate from the domain event topics.
const (
	TopicDocumentProcess      = "masterfabric.jobs.document-process"
	TopicNotificationDispatch = "masterfabric.jobs.notification-dispatch"
	TopicEmailSend            = "masterfabric.jobs.email-send"
	TopicDeadLetter           = "masterfabric.jobs.dead-letter"
)

// AllTopics returns every job topic (used to ensure Kafka topics at startup).
func AllTopics() []string {
	return []string{
		TopicDocumentProcess,
		TopicNotificationDispatch,
		TopicEmailSend,
		TopicDeadLetter,
	}
}

// Enqueuer publishes a job for asynchronous background processing. It is a thin
// abstraction over the event bus so handlers/use-cases stay transport-agnostic.
type Enqueuer interface {
	Enqueue(ctx context.Context, topic string, job interface{}) error
}

// DocumentProcessJob is queued after an evrak upload to run post-processing
// (checksum verification, virus scan, indexing) off the request path.
type DocumentProcessJob struct {
	DocumentID     string `json:"document_id"`
	OrganizationID string `json:"organization_id"`
	AppID          string `json:"app_id"`
	StorageBucket  string `json:"storage_bucket"`
	StorageKey     string `json:"storage_key"`
	Checksum       string `json:"checksum"`
}

// NotificationDispatchJob is queued to deliver a notification asynchronously.
type NotificationDispatchJob struct {
	NotificationID string `json:"notification_id"`
	OrganizationID string `json:"organization_id"`
	AppID          string `json:"app_id"`
	Channel        string `json:"channel"`
}

// EmailSendJob is queued to send an email asynchronously.
type EmailSendJob struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// DeadLetter wraps a job that exhausted its retries for later inspection.
type DeadLetter struct {
	OriginalTopic string `json:"original_topic"`
	Payload       string `json:"payload"`
	Error         string `json:"error"`
	FailedAt      string `json:"failed_at"`
}
