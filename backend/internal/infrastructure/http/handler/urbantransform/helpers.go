package urbantransform

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

// queryInt parses an integer query parameter, returning 0 when absent or invalid.
func queryInt(r *http.Request, key string) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return 0
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return n
}

// parseUUIDQuery parses a UUID query parameter, returning nil when absent or invalid.
func parseUUIDQuery(r *http.Request, key string) *uuid.UUID {
	val := r.URL.Query().Get(key)
	if val == "" {
		return nil
	}
	id, err := uuid.Parse(val)
	if err != nil {
		return nil
	}
	return &id
}
