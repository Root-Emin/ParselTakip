package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration.
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	Kafka     KafkaConfig
	Log       LogConfig
	Storage   StorageConfig
	Security  SecurityConfig
	SeedAdmin SeedAdminConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int32
	MinConns int32
}

// DSN returns the PostgreSQL connection string.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Addr returns the Redis address string.
func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// JWTConfig holds JWT signing settings.
type JWTConfig struct {
	Secret          string
	ExpirationHours int
	Issuer          string
}

// KafkaConfig holds Kafka connection and consumer settings.
type KafkaConfig struct {
	Brokers           []string
	GroupID           string
	Enabled           bool
	NumPartitions     int
	ReplicationFactor int
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level  string // debug, info, warn, error
	Format string // json, text
}

// StorageConfig holds MinIO/S3-compatible object storage settings used for
// document (evrak) uploads.
type StorageConfig struct {
	Enabled        bool
	Endpoint       string        // internal endpoint used by the server (e.g. localhost:9000)
	PublicEndpoint string        // endpoint embedded in presigned URLs handed to clients; defaults to Endpoint
	AccessKey      string
	SecretKey      string
	Bucket         string
	Region         string
	UseSSL         bool
	PresignExpiry  time.Duration // TTL for presigned GET/PUT URLs (KVKK: keep short)
}

// SecurityConfig holds security and KVKK-related settings.
type SecurityConfig struct {
	PIIEncryptionKey string        // secret used to derive AES-256-GCM key for PII at rest
	BcryptCost       int           // bcrypt cost factor for password hashing
	MaxLoginAttempts int           // failed attempts before account lockout
	LockoutDuration  time.Duration // how long an account stays locked
}

// SeedAdminConfig holds the bootstrap super-admin that is ensured on startup.
type SeedAdminConfig struct {
	Enabled  bool
	Email    string
	Password string
	OrgName  string
	OrgSlug  string
	AppName  string
	AppSlug  string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         envOrDefault("SERVER_HOST", "0.0.0.0"),
			Port:         envOrDefaultInt("SERVER_PORT", 8080),
			ReadTimeout:  time.Duration(envOrDefaultInt("SERVER_READ_TIMEOUT_SECONDS", 15)) * time.Second,
			WriteTimeout: time.Duration(envOrDefaultInt("SERVER_WRITE_TIMEOUT_SECONDS", 15)) * time.Second,
			IdleTimeout:  time.Duration(envOrDefaultInt("SERVER_IDLE_TIMEOUT_SECONDS", 60)) * time.Second,
		},
		Database: DatabaseConfig{
			Host:     envOrDefault("DB_HOST", "localhost"),
			Port:     envOrDefaultInt("DB_PORT", 5432),
			User:     envOrDefault("DB_USER", "masterfabric"),
			Password: envOrDefault("DB_PASSWORD", "masterfabric"),
			DBName:   envOrDefault("DB_NAME", "masterfabric"),
			SSLMode:  envOrDefault("DB_SSLMODE", "disable"),
			MaxConns: int32(envOrDefaultInt("DB_MAX_CONNS", 25)),
			MinConns: int32(envOrDefaultInt("DB_MIN_CONNS", 5)),
		},
		Redis: RedisConfig{
			Host:     envOrDefault("REDIS_HOST", "localhost"),
			Port:     envOrDefaultInt("REDIS_PORT", 6379),
			Password: envOrDefault("REDIS_PASSWORD", ""),
			DB:       envOrDefaultInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:          envOrDefault("JWT_SECRET", "change-me-in-production"),
			ExpirationHours: envOrDefaultInt("JWT_EXPIRATION_HOURS", 24),
			Issuer:          envOrDefault("JWT_ISSUER", "masterfabric"),
		},
		Kafka: KafkaConfig{
			Brokers:           envOrDefaultSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			GroupID:           envOrDefault("KAFKA_GROUP_ID", "masterfabric-go"),
			Enabled:           envOrDefault("KAFKA_ENABLED", "false") == "true",
			NumPartitions:     envOrDefaultInt("KAFKA_NUM_PARTITIONS", 3),
			ReplicationFactor: envOrDefaultInt("KAFKA_REPLICATION_FACTOR", 1),
		},
		Log: LogConfig{
			Level:  envOrDefault("LOG_LEVEL", "info"),
			Format: envOrDefault("LOG_FORMAT", "json"),
		},
		Storage: StorageConfig{
			Enabled:        envOrDefault("MINIO_ENABLED", "true") == "true",
			Endpoint:       envOrDefault("MINIO_ENDPOINT", "localhost:9000"),
			PublicEndpoint: envOrDefault("MINIO_PUBLIC_ENDPOINT", ""),
			AccessKey:      envOrDefault("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey:      envOrDefault("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:         envOrDefault("MINIO_BUCKET", "parseltakip"),
			Region:         envOrDefault("MINIO_REGION", "us-east-1"),
			UseSSL:         envOrDefault("MINIO_USE_SSL", "false") == "true",
			PresignExpiry:  time.Duration(envOrDefaultInt("MINIO_PRESIGN_EXPIRY_MINUTES", 10)) * time.Minute,
		},
		Security: SecurityConfig{
			PIIEncryptionKey: envOrDefault("PII_ENCRYPTION_KEY", "dev-pii-key-change-me-in-production"),
			BcryptCost:       envOrDefaultInt("BCRYPT_COST", 12),
			MaxLoginAttempts: envOrDefaultInt("MAX_LOGIN_ATTEMPTS", 5),
			LockoutDuration:  time.Duration(envOrDefaultInt("LOGIN_LOCKOUT_MINUTES", 15)) * time.Minute,
		},
		SeedAdmin: SeedAdminConfig{
			Enabled:  envOrDefault("SEED_ADMIN_ENABLED", "true") == "true",
			Email:    envOrDefault("SEED_ADMIN_EMAIL", "muhammedysnvurucu@gmail.com"),
			Password: envOrDefault("SEED_ADMIN_PASSWORD", "cursor123"),
			OrgName:  envOrDefault("SEED_ADMIN_ORG_NAME", "ParselTakip"),
			OrgSlug:  envOrDefault("SEED_ADMIN_ORG_SLUG", "parseltakip"),
			AppName:  envOrDefault("SEED_ADMIN_APP_NAME", "ParselTakip"),
			AppSlug:  envOrDefault("SEED_ADMIN_APP_SLUG", "parseltakip"),
		},
	}
}

func envOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func envOrDefaultInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func envOrDefaultSlice(key string, defaultVal []string) []string {
	if val := os.Getenv(key); val != "" {
		parts := strings.Split(val, ",")
		var result []string
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultVal
}
