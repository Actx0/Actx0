// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package migration

import (
	"database/sql"
	"fmt"
	"strings"
)

// detectDriver attempts to determine the database driver type
func detectDriver(db *sql.DB) string {
	_, err := db.Exec("SELECT version()")
	if err == nil {
		var version string
		db.QueryRow("SELECT version()").Scan(&version)
		if strings.Contains(strings.ToLower(version), "postgresql") {
			return "postgres"
		}
	}

	// Unknown database driver
	return "unknown"
}

// GetAll returns all registered migrations (basic tables + users + sessions only)
func GetAll() []Migration {
	return []Migration{
		{
			Version:     "20250101000001",
			Description: "Create configs table",
			Up:          createConfigsTable,
			Down:        dropConfigsTable,
		},
		{
			Version:     "20250101000002",
			Description: "Create users table",
			Up:          createUsersTable,
			Down:        dropUsersTable,
		},
		{
			Version:     "20250101000003",
			Description: "Create users_meta table",
			Up:          createUsersMetaTable,
			Down:        dropUsersMetaTable,
		},
		{
			Version:     "20250101000004",
			Description: "Create user_sessions table",
			Up:          createUserSessionsTable,
			Down:        dropUserSessionsTable,
		}, {
			Version:     "20250101000005",
			Description: "Create user_api_keys table",
			Up:          createUserAPIKeysTable,
			Down:        dropUserAPIKeysTable,
		},
		{
			Version:     "20250101000006",
			Description: "Create workspaces table",
			Up:          createWorkspacesTable,
			Down:        dropWorkspacesTable,
		},
		{
			Version:     "20250101000007",
			Description: "Create user_invites table",
			Up:          createUserInvitesTable,
			Down:        dropUserInvitesTable,
		},
		{
			Version:     "20250101000008",
			Description: "Create workspaces_meta table",
			Up:          createWorkspacesMetaTable,
			Down:        dropWorkspacesMetaTable,
		},
		{
			Version:     "20250101000009",
			Description: "Create workspace_users table",
			Up:          createWorkspaceUsersTable,
			Down:        dropWorkspaceUsersTable,
		},
	}
}

// createConfigsTable creates the configs table
func createConfigsTable(db *sql.DB) error {
	driver := detectDriver(db)
	var query string

	switch driver {
	case "postgres":
		query = `
		CREATE TABLE configs (
			id UUID PRIMARY KEY,
			key VARCHAR(128) NOT NULL UNIQUE,
			value TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX idx_configs_key ON configs(key)`
	default:
		return fmt.Errorf("unsupported database driver: %s", driver)
	}

	_, err := db.Exec(query)
	return err
}

// dropConfigsTable drops the configs table
func dropConfigsTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS configs")
	return err
}

// createUsersTable creates the users table
func createUsersTable(db *sql.DB) error {
	driver := detectDriver(db)
	var query string

	switch driver {
	case "postgres":
		query = `
		CREATE TABLE users (
			id UUID PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(191) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'regular',
			is_active BOOLEAN DEFAULT true,
			last_login_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX idx_email ON users(email)`
	default:
		return fmt.Errorf("unsupported database driver: %s", driver)
	}

	_, err := db.Exec(query)
	return err
}

// dropUsersTable drops the users table
func dropUsersTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	return err
}

// createUsersMetaTable creates the users_meta table
func createUsersMetaTable(db *sql.DB) error {
	driver := detectDriver(db)
	var query string

	switch driver {
	case "postgres":
		query = `
		CREATE TABLE users_meta (
			id UUID PRIMARY KEY,
			key VARCHAR(128) NOT NULL,
			value TEXT,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (user_id, key)
		);
		CREATE INDEX idx_user_id ON users_meta(user_id);
		CREATE INDEX idx_key ON users_meta(key)`
	default:
		return fmt.Errorf("unsupported database driver: %s", driver)
	}

	_, err := db.Exec(query)
	return err
}

// dropUsersMetaTable drops the users_meta table
func dropUsersMetaTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users_meta")
	return err
}

// createUserSessionsTable creates the user_sessions table
func createUserSessionsTable(db *sql.DB) error {
	driver := detectDriver(db)
	var query string

	switch driver {
	case "postgres":
		query = `
		CREATE TABLE user_sessions (
			id UUID PRIMARY KEY,
			token VARCHAR(255) NOT NULL UNIQUE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			ip_address VARCHAR(45),
			user_agent VARCHAR(255),
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX idx_token ON user_sessions(token);
		CREATE INDEX idx_user_id ON user_sessions(user_id);
		CREATE INDEX idx_expires_at ON user_sessions(expires_at)`
	default:
		return fmt.Errorf("unsupported database driver: %s", driver)
	}

	_, err := db.Exec(query)
	return err
}

// dropUserSessionsTable drops the user_sessions table
func dropUserSessionsTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS user_sessions")
	return err
}

// createUserAPIKeysTable creates the user_api_keys table
func createUserAPIKeysTable(db *sql.DB) error {
	driver := detectDriver(db)
	var query string

	switch driver {
	case "postgres":
		query = `
		CREATE TABLE user_api_keys (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(100),
			key VARCHAR(512) NOT NULL,
			expires_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX idx_user_api_keys_user_id ON user_api_keys(user_id);
		CREATE INDEX idx_user_api_keys_expires_at ON user_api_keys(expires_at);
		CREATE INDEX idx_user_api_keys_key ON user_api_keys(key)`
	default:
		return fmt.Errorf("unsupported database driver: %s", driver)
	}

	_, err := db.Exec(query)
	return err
}

// dropUserAPIKeysTable drops the user_api_keys table
func dropUserAPIKeysTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS user_api_keys")
	return err
}

// createUserInvitesTable creates the user_invites table
func createUserInvitesTable(db *sql.DB) error {
	driver := detectDriver(db)
	var query string

	switch driver {
	case "postgres":
		query = `
		CREATE TABLE user_invites (
			id UUID PRIMARY KEY,
			email VARCHAR(191) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'regular',
			token VARCHAR(255) NOT NULL UNIQUE,
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			inviter_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			expires_at TIMESTAMP NOT NULL,
			accepted_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX idx_user_invites_email ON user_invites(email);
		CREATE INDEX idx_user_invites_token ON user_invites(token);
		CREATE INDEX idx_user_invites_expires_at ON user_invites(expires_at)`
	default:
		return fmt.Errorf("unsupported database driver: %s", driver)
	}

	_, err := db.Exec(query)
	return err
}

// dropUserInvitesTable drops the user_invites table.
func dropUserInvitesTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS user_invites")
	return err
}

// createWorkspacesTable creates the workspaces table
func createWorkspacesTable(db *sql.DB) error {
	driver := detectDriver(db)
	var query string

	switch driver {
	case "postgres":
		query = `
		CREATE TABLE workspaces (
			id UUID PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
	default:
		return fmt.Errorf("unsupported database driver: %s", driver)
	}

	_, err := db.Exec(query)
	return err
}

// dropWorkspacesTable drops the workspaces table.
func dropWorkspacesTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS workspaces")
	return err
}

// createWorkspacesMetaTable creates the workspaces_meta table (key-value per workspace)
func createWorkspacesMetaTable(db *sql.DB) error {
	driver := detectDriver(db)
	var query string

	switch driver {
	case "postgres":
		query = `
		CREATE TABLE workspaces_meta (
			id UUID PRIMARY KEY,
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			key VARCHAR(128) NOT NULL,
			value TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (workspace_id, key)
		);
		CREATE INDEX idx_workspaces_meta_workspace_id ON workspaces_meta(workspace_id);
		CREATE INDEX idx_workspaces_meta_key ON workspaces_meta(key)`
	default:
		return fmt.Errorf("unsupported database driver: %s", driver)
	}

	_, err := db.Exec(query)
	return err
}

// dropWorkspacesMetaTable drops the workspaces_meta table.
func dropWorkspacesMetaTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS workspaces_meta")
	return err
}

// createWorkspaceUsersTable creates the workspace_users table (user-workspace membership and role)
func createWorkspaceUsersTable(db *sql.DB) error {
	driver := detectDriver(db)
	var query string

	switch driver {
	case "postgres":
		query = `
		CREATE TABLE workspace_users (
			id UUID PRIMARY KEY,
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			role VARCHAR(50) NOT NULL DEFAULT 'regular',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (workspace_id, user_id)
		);
		CREATE INDEX idx_workspace_users_workspace_id ON workspace_users(workspace_id);
		CREATE INDEX idx_workspace_users_user_id ON workspace_users(user_id)`
	default:
		return fmt.Errorf("unsupported database driver: %s", driver)
	}

	_, err := db.Exec(query)
	return err
}

// dropWorkspaceUsersTable drops the workspace_users table.
func dropWorkspaceUsersTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS workspace_users")
	return err
}
