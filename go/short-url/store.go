package main

import (
	"fmt"
	"sync"
)

// URLStore defines the interface for storing and retrieving URL mappings.
type URLStore interface {
	Set(id, url string) error
	Get(id string) (string, bool)
	GenerateID() string // Method to generate a new short ID
}

// InMemoryURLStore implements URLStore using an in-memory map.
type InMemoryURLStore struct {
	mu    sync.RWMutex
	urls  map[string]string
	count int // Simple counter for generating IDs
}

// NewInMemoryURLStore creates a new InMemoryURLStore.
func NewInMemoryURLStore() *InMemoryURLStore {
	return &InMemoryURLStore{
		urls: make(map[string]string),
	}
}

// Set stores the mapping between a short ID and a long URL.
func (s *InMemoryURLStore) Set(id, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urls[id] = url
	return nil
}

// Get retrieves the long URL associated with a short ID.
func (s *InMemoryURLStore) Get(id string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.urls[id]
	return url, ok
}

// GenerateID generates a simple sequential short ID.
// TODO: Implement a more robust ID generation strategy (e.g., random strings, base62 encoding).
func (s *InMemoryURLStore) GenerateID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count++
	// Simple base-10 representation for now
	return fmt.Sprintf("%d", s.count)
}
