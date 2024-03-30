package store

import "time"

type Store interface {
	Get(key string) (int, error)
	IsBlocked(key string) (bool, error)
	Increment(key string, expiration time.Duration) (int, error)
	Block(key string, blockDuration time.Duration) error
}
