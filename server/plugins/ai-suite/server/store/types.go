package store

import "errors"

var (
	// ErrNotFound is returned when a key does not exist.
	ErrNotFound = errors.New("store: not found")

	// ErrInvalidKey indicates a missing or malformed key value.
	ErrInvalidKey = errors.New("store: invalid key")
)
