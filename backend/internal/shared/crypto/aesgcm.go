// Package crypto provides authenticated symmetric encryption (AES-256-GCM)
// for protecting sensitive PII fields at rest (KVKK compliance), plus a keyed
// blind-index (HMAC-SHA256) for equality lookups on encrypted columns.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// Encryptor encrypts/decrypts strings with AES-256-GCM and derives blind indexes.
type Encryptor struct {
	gcm     cipher.AEAD
	hmacKey []byte
}

// NewEncryptor derives a 32-byte AES key from the provided secret via SHA-256,
// so any non-empty secret string is accepted. The HMAC key for blind indexes is
// derived from the same secret using a different domain separator.
func NewEncryptor(secret string) (*Encryptor, error) {
	if secret == "" {
		return nil, errors.New("encryption secret must not be empty")
	}
	aesKey := sha256.Sum256([]byte("parseltakip:aes:" + secret))
	hmacKey := sha256.Sum256([]byte("parseltakip:hmac:" + secret))

	block, err := aes.NewCipher(aesKey[:])
	if err != nil {
		return nil, fmt.Errorf("aes new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}
	return &Encryptor{gcm: gcm, hmacKey: hmacKey[:]}, nil
}

// Encrypt returns base64(nonce || ciphertext) for plaintext. Empty input yields
// empty output so optional fields round-trip cleanly.
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("read nonce: %w", err)
	}
	sealed := e.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt reverses Encrypt. Empty input yields empty output.
func (e *Encryptor) Decrypt(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	nonceSize := e.gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := e.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("gcm open: %w", err)
	}
	return string(plaintext), nil
}

// BlindIndex returns a deterministic hex HMAC-SHA256 of value, enabling equality
// search on an encrypted column without storing or revealing the plaintext.
func (e *Encryptor) BlindIndex(value string) string {
	if value == "" {
		return ""
	}
	mac := hmac.New(sha256.New, e.hmacKey)
	mac.Write([]byte(value))
	return hex.EncodeToString(mac.Sum(nil))
}
