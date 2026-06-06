// Package storage provides a MinIO/S3-compatible implementation of the
// domain storage.ObjectStorage interface for document uploads.
package storage

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	domainStorage "github.com/masterfabric-go/masterfabric/internal/domain/storage"
	"github.com/masterfabric-go/masterfabric/internal/shared/config"
)

// MinIOStorage implements domain storage.ObjectStorage backed by MinIO.
type MinIOStorage struct {
	client         *minio.Client
	bucket         string
	publicEndpoint string
	useSSL         bool
	logger         *slog.Logger
}

// Verify interface compliance at compile time.
var _ domainStorage.ObjectStorage = (*MinIOStorage)(nil)

// NewMinIOStorage connects to MinIO and ensures the configured bucket exists.
func NewMinIOStorage(ctx context.Context, cfg config.StorageConfig, logger *slog.Logger) (*MinIOStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("minio new client: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("minio bucket exists: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{Region: cfg.Region}); err != nil {
			return nil, fmt.Errorf("minio make bucket: %w", err)
		}
		logger.Info("minio bucket created", "bucket", cfg.Bucket)
	}

	return &MinIOStorage{
		client:         client,
		bucket:         cfg.Bucket,
		publicEndpoint: cfg.PublicEndpoint,
		useSSL:         cfg.UseSSL,
		logger:         logger,
	}, nil
}

// Bucket returns the configured bucket name.
func (s *MinIOStorage) Bucket() string { return s.bucket }

// Upload stores an object and returns its metadata.
func (s *MinIOStorage) Upload(ctx context.Context, in domainStorage.UploadInput) (*domainStorage.ObjectInfo, error) {
	info, err := s.client.PutObject(ctx, s.bucket, in.Key, in.Reader, in.Size, minio.PutObjectOptions{
		ContentType: in.ContentType,
	})
	if err != nil {
		return nil, fmt.Errorf("minio put object: %w", err)
	}
	return &domainStorage.ObjectInfo{
		Key:         in.Key,
		Size:        info.Size,
		ContentType: in.ContentType,
		ETag:        info.ETag,
	}, nil
}

// PresignGet returns a short-lived download URL.
func (s *MinIOStorage) PresignGet(ctx context.Context, key string, expiry time.Duration, downloadName string) (string, error) {
	reqParams := make(url.Values)
	if downloadName != "" {
		reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=%q", downloadName))
	}
	u, err := s.client.PresignedGetObject(ctx, s.bucket, key, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("minio presign get: %w", err)
	}
	return s.rewriteHost(u), nil
}

// PresignPut returns a short-lived upload URL for direct client PUTs.
func (s *MinIOStorage) PresignPut(ctx context.Context, key string, expiry time.Duration) (string, error) {
	u, err := s.client.PresignedPutObject(ctx, s.bucket, key, expiry)
	if err != nil {
		return "", fmt.Errorf("minio presign put: %w", err)
	}
	return s.rewriteHost(u), nil
}

// Delete removes an object.
func (s *MinIOStorage) Delete(ctx context.Context, key string) error {
	if err := s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("minio remove object: %w", err)
	}
	return nil
}

// Stat returns metadata for an object.
func (s *MinIOStorage) Stat(ctx context.Context, key string) (*domainStorage.ObjectInfo, error) {
	info, err := s.client.StatObject(ctx, s.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("minio stat object: %w", err)
	}
	return &domainStorage.ObjectInfo{
		Key:         key,
		Size:        info.Size,
		ContentType: info.ContentType,
		ETag:        info.ETag,
	}, nil
}

// rewriteHost swaps the signed URL host for the configured public endpoint so
// browsers can reach MinIO even when the server uses an internal hostname.
func (s *MinIOStorage) rewriteHost(u *url.URL) string {
	if s.publicEndpoint == "" {
		return u.String()
	}
	public := s.publicEndpoint
	if !strings.Contains(public, "://") {
		scheme := "http"
		if s.useSSL {
			scheme = "https"
		}
		public = scheme + "://" + public
	}
	pub, err := url.Parse(public)
	if err != nil {
		return u.String()
	}
	u.Scheme = pub.Scheme
	u.Host = pub.Host
	return u.String()
}
