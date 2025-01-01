// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/clivern/actx0/pkg/util"
)

type localStore struct {
	config Config
}

// newLocalStore creates a filesystem-backed document store.
func newLocalStore(config Config) (*localStore, error) {
	err := util.EnsureDir(config.Local.Path, 0o755)
	if err != nil {
		return nil, fmt.Errorf("local storage ensure dir: %w", err)
	}

	return &localStore{config: config}, nil
}

// Config returns the store configuration.
func (s *localStore) Config() Config {
	return s.config
}

// PutDocument writes document bytes to storage.
func (s *localStore) PutDocument(_ context.Context, id string, text string) (string, error) {
	key := DocumentKey(id)

	path := s.objectPath(key)
	err := util.EnsureDir(filepath.Dir(path), 0o755)
	if err != nil {
		return "", fmt.Errorf("local storage ensure dir: %w", err)
	}

	err = os.WriteFile(path, []byte(text), 0o644)
	if err != nil {
		return "", fmt.Errorf("local storage put document: %w", err)
	}

	return s.documentURI(key)
}

// GetDocument reads document bytes from storage.
func (s *localStore) GetDocument(_ context.Context, id string) (string, error) {
	key := DocumentKey(id)

	data, err := os.ReadFile(s.objectPath(key))
	if os.IsNotExist(err) {
		return "", ErrDocumentNotFound
	}
	if err != nil {
		return "", fmt.Errorf("local storage get document: %w", err)
	}

	return string(data), nil
}

// DeleteDocument removes a document from storage.
func (s *localStore) DeleteDocument(_ context.Context, id string) error {
	key := DocumentKey(id)

	err := os.Remove(s.objectPath(key))
	if os.IsNotExist(err) {
		return ErrDocumentNotFound
	}
	if err != nil {
		return fmt.Errorf("local storage delete document: %w", err)
	}

	return nil
}

// DocumentURI returns the URI for a stored document.
func (s *localStore) DocumentURI(id string) (string, error) {
	return s.documentURI(DocumentKey(id))
}

// objectPath returns the filesystem path for a document key.
func (s *localStore) objectPath(key string) string {
	return filepath.Join(s.config.Local.Path, filepath.FromSlash(key))
}

// documentURI builds the URI for a stored document.
func (s *localStore) documentURI(key string) (string, error) {
	path := s.objectPath(key)
	return "file://" + filepath.ToSlash(path), nil
}
