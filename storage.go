package main

import (
	"sync"
)

// QAPair represents a question and answer pair
type QAPair struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// Storage interface defines methods for managing Q&A history
type Storage interface {
	Add(pair QAPair) error
	GetAll() ([]QAPair, error)
	Close() error
}

// MemoryStorage implements Storage interface using in-memory storage
type MemoryStorage struct {
	maxSize int
	mu      sync.RWMutex
	pairs   []QAPair
}

// NewMemoryStorage creates a new MemoryStorage instance
func NewMemoryStorage(maxSize int) *MemoryStorage {
	return &MemoryStorage{
		maxSize: maxSize,
		pairs:   make([]QAPair, 0),
	}
}

// Add adds a new Q&A pair to storage
func (s *MemoryStorage) Add(pair QAPair) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pairs = append(s.pairs, pair)

	// Enforce maximum size by removing oldest entries
	if len(s.pairs) > s.maxSize {
		s.pairs = s.pairs[len(s.pairs)-s.maxSize:]
	}

	return nil
}

// GetAll returns all Q&A pairs
func (s *MemoryStorage) GetAll() ([]QAPair, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make([]QAPair, len(s.pairs))
	copy(result, s.pairs)
	return result, nil
}

// Close performs cleanup
func (s *MemoryStorage) Close() error {
	// No resources to clean up for in-memory storage
	return nil
}
