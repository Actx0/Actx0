// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"database/sql"
)

const UsageTypeAPICalls = "api_calls"
const UsageTypeAITokens = "ai_tokens"
const UsageTypeAICost = "ai_cost"
const UsageUnitCalls = "calls"
const UsageUnitTokens = "tokens"
const UsageUnitNanoUSD = "nano_usd"

// WorkspaceStats holds aggregate metrics for a workspace dashboard.
type WorkspaceStats struct {
	MemoriesStored int64
	APICallsMonth  int64
	ActiveAgents   int64
	Prompts        int64
	TopAgents      []*TopAgentStats
}

// TopAgentStats is an agent ranked by session and memory counts.
type TopAgentStats struct {
	Id       Id
	Name     string
	Handle   string
	Sessions int64
	Memories int64
}

// WorkspaceStatsRepository loads workspace dashboard metrics.
type WorkspaceStatsRepository interface {
	GetByWorkspaceId(workspaceId Id) (*WorkspaceStats, error)
}

type WorkspaceStatsRepositoryPostgres struct {
	db *sql.DB
}

// NewWorkspaceStatsRepository returns a workspace stats repository.
func NewWorkspaceStatsRepository(db *sql.DB) WorkspaceStatsRepository {
	return &WorkspaceStatsRepositoryPostgres{db: db}
}

// GetByWorkspaceId returns a workspace stats by workspace id.
func (r *WorkspaceStatsRepositoryPostgres) GetByWorkspaceId(workspaceId Id) (*WorkspaceStats, error) {
	stats := &WorkspaceStats{}

	// Memories Stored
	err := r.db.QueryRow(
		`SELECT COUNT(*)
		FROM session_memories sm
		INNER JOIN sessions s ON s.id = sm.session_id
		WHERE s.workspace_id = $1`,
		workspaceId.String(),
	).Scan(&stats.MemoriesStored)
	if err != nil {
		return nil, err
	}

	// APICalls for the current subscription period
	err = r.db.QueryRow(
		`SELECT COALESCE(SUM(quantity), 0)
		FROM usage
		WHERE workspace_id = $1
			AND type = $2
			AND period_start = (SELECT period_start FROM subscriptions WHERE workspace_id = $1)`,
		workspaceId.String(),
		UsageTypeAPICalls,
	).Scan(&stats.APICallsMonth)
	if err != nil {
		return nil, err
	}

	// Active Agents
	err = r.db.QueryRow(
		`SELECT COUNT(*)
		FROM agents
		WHERE workspace_id = $1 AND status = 'active'`,
		workspaceId.String(),
	).Scan(&stats.ActiveAgents)
	if err != nil {
		return nil, err
	}

	// Prompts
	err = r.db.QueryRow(
		`SELECT COUNT(*)
		FROM prompts
		WHERE workspace_id = $1`,
		workspaceId.String(),
	).Scan(&stats.Prompts)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(
		`SELECT
			a.id,
			a.name,
			a.handle,
			COALESCE(sc.sessions_count, 0) AS sessions_count,
			COALESCE(mc.memories_count, 0) AS memories_count
		FROM agents a
		LEFT JOIN (
			SELECT agent_id, COUNT(*) AS sessions_count
			FROM sessions
			WHERE workspace_id = $1
			GROUP BY agent_id
		) sc ON sc.agent_id = a.id
		LEFT JOIN (
			SELECT s.agent_id, COUNT(*) AS memories_count
			FROM session_memories sm
			INNER JOIN sessions s ON s.id = sm.session_id
			WHERE s.workspace_id = $1
			GROUP BY s.agent_id
		) mc ON mc.agent_id = a.id
		WHERE a.workspace_id = $1
			AND (COALESCE(sc.sessions_count, 0) > 0 OR COALESCE(mc.memories_count, 0) > 0)
		ORDER BY sessions_count DESC, memories_count DESC, a.name ASC
		LIMIT 5`,
		workspaceId.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats.TopAgents = []*TopAgentStats{}
	for rows.Next() {
		agent := &TopAgentStats{}
		err = rows.Scan(
			&agent.Id,
			&agent.Name,
			&agent.Handle,
			&agent.Sessions,
			&agent.Memories,
		)
		if err != nil {
			return nil, err
		}
		stats.TopAgents = append(stats.TopAgents, agent)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
