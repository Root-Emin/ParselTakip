// Package storage defines the object-storage abstraction used for document
// (evrak) uploads. Implementations live in infrastructure (e.g. MinIO/S3).
package storage

import (
	"context"
	"io"
	"time"
)

// UploadInput carries the data needed to persist an object.
type UploadInput struct {
	Key         string
	Reader      io.Reader
	Size        int64 // -1 when unknown (streamed)
	ContentType string
}

// ObjectInfo describes a stored object.
type ObjectInfo struct {
	Key         string
	Size        int64
	ContentType string
	ETag        string
}

// ObjectStorage abstracts an S3-compatible object store.
type ObjectStorage interface {
	// Upload stores an object and returns its metadata.
	Upload(ctx context.Context, in UploadInput) (*ObjectInfo, error)
	// PresignGet returns a short-lived URL to download an object. downloadName,
	// when set, forces a Content-Disposition attachment filename.
	PresignGet(ctx context.Context, key string, expiry time.Duration, downloadName string) (string, error)
	// PresignPut returns a short-lived URL the client can PUT an object to directly.
	PresignPut(ctx context.Context, key string, expiry time.Duration) (string, error)
	// Delete removes an object.
	Delete(ctx context.Context, key string) error
	// Stat returns metadata for an object.
	Stat(ctx context.Context, key string) (*ObjectInfo, error)
	// Bucket returns the configured bucket name.
	Bucket() string
}
