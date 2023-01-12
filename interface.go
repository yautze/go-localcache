package localcache

import "errors"

var (
	// ErrDataNotFound means data not found by this key
	ErrDataNotFound = errors.New("data not found by this key")
)

// Cache porvide operate methods
type Cache interface {
	// Set value into cache with key
	Set(key string, val interface{}) error

	// Get value form cache by key
	Get(key string) (interface{}, error)
}
