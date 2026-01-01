// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	// dbConn holds the singleton database connection
	dbConn *Connection
	// mu protects dbConn during initialization
	mu sync.RWMutex
)

// InitDB initializes the global database connection
func InitDB(config Config) error {
	mu.Lock()
	defer mu.Unlock()

	if dbConn != nil {
		log.Warn().Msg("Database connection already initialized")
		return nil
	}

	conn, err := NewConnection(config)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	dbConn = conn
	log.Info().Msg("Global database connection initialized")
	return nil
}

// GetDB returns the global database connection
func GetDB() *sql.DB {
	mu.RLock()
	defer mu.RUnlock()

	if dbConn == nil {
		log.Error().Msg("Database not initialized")
		return nil
	}

	return dbConn.DB
}

// CloseDB closes the global database connection
func CloseDB() error {
	mu.Lock()
	defer mu.Unlock()

	if dbConn == nil {
		return nil
	}

	err := dbConn.Close()
	dbConn = nil
	return err
}
