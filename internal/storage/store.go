package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	path string
	mu   sync.Mutex
}

func New(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("storage path cannot be empty")
	}
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create storage directory: %w", err)
		}
	}
	return &Store{path: path}, nil
}

func (s *Store) Load(v interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open storage file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("failed to decode storage data: %w", err)
	}
	return nil
}

func (s *Store) Save(v interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tempFile := s.path + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temp storage file: %w", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		file.Close()
		os.Remove(tempFile)
		return fmt.Errorf("failed to encode storage data: %w", err)
	}
	file.Close()

	if err := os.Rename(tempFile, s.path); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to replace storage file: %w", err)
	}
	return nil
}
