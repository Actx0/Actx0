// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"github.com/clivern/actx0/db"
)

// Stats is the module for workspace dashboard metrics.
type Stats struct {
	WorkspaceRepository db.WorkspaceRepository
	StatsRepository     db.WorkspaceStatsRepository
}

// NewStats creates a stats module with the given repositories.
func NewStats(workspaces db.WorkspaceRepository, stats db.WorkspaceStatsRepository) *Stats {
	return &Stats{
		WorkspaceRepository: workspaces,
		StatsRepository:     stats,
	}
}

// WorkspaceStatsResponse is workspace metrics shaped for API responses.
type WorkspaceStatsResponse struct {
	MemoriesStored int64               `json:"memoriesStored"`
	APICallsMonth  int64               `json:"apiCallsMonth"`
	ActiveAgents   int64               `json:"activeAgents"`
	Prompts        int64               `json:"prompts"`
	TopAgents      []*TopAgentResponse `json:"topAgents"`
}

// TopAgentResponse is a top agent ranked by activity.
type TopAgentResponse struct {
	Id       db.Id  `json:"id"`
	Name     string `json:"name"`
	Handle   string `json:"handle"`
	Sessions int64  `json:"sessions"`
	Memories int64  `json:"memories"`
}

// GetWorkspaceStats returns dashboard metrics for a workspace.
func (s *Stats) GetWorkspaceStats(workspaceId db.Id) (*WorkspaceStatsResponse, error) {
	workspace, err := s.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	stats, err := s.StatsRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}

	topAgents := make([]*TopAgentResponse, 0, len(stats.TopAgents))
	for _, agent := range stats.TopAgents {
		topAgents = append(topAgents, &TopAgentResponse{
			Id:       agent.Id,
			Name:     agent.Name,
			Handle:   agent.Handle,
			Sessions: agent.Sessions,
			Memories: agent.Memories,
		})
	}

	return &WorkspaceStatsResponse{
		MemoriesStored: stats.MemoriesStored,
		APICallsMonth:  stats.APICallsMonth,
		ActiveAgents:   stats.ActiveAgents,
		Prompts:        stats.Prompts,
		TopAgents:      topAgents,
	}, nil
}
