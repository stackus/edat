package msg

import (
	"fmt"
)

// Headers a map of strings keyed by Message header keys
type Headers map[string]string

// Has returned whether or not the given key exists in the headers
func (h Headers) Has(key string) bool {
	_, exists := h[key]

	return exists
}

// Get returns the value for the given key. Returns a blank string if it does not exist
func (h Headers) Get(key string) string {
	return h[key]
}

// GetRequired returns the value for the given key. Returns an error if it does not exist
func (h Headers) GetRequired(key string) (string, error) {
	value, exists := h[key]
	if !exists {
		return "", fmt.Errorf("missing required header `%s`", key)
	}

	return value, nil
}

// Set sets or overwrites the key with the value
func (h Headers) Set(key, value string) {
	h[key] = value
}
