package data

import (
	"errors"
	"sync"
)

var (
	// SnapshotFileName is the name of the file that will be created when
	SnapshotFileName = "snapshot.json"

	// ErrKeyNotFound is the error message when a key is not found
	ErrKeyNotFound = errors.New("key not found")
)

type Store struct {
	KeyValue         map[string]string
	mu               sync.RWMutex
	snapShotLocation string
}

func NewStore(snapShotLocation string) *Store {
	return &Store{
		KeyValue:         make(map[string]string),
		mu:               sync.RWMutex{},
		snapShotLocation: snapShotLocation,
	}
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.KeyValue[key] = value
}

func (s *Store) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.KeyValue[key]
	if !ok {
		return "", ErrKeyNotFound
	}

	return value, nil
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.KeyValue, key)
}
