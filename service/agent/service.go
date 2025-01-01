// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.
//
// Agent service: session memory, indexing, and semantic search.

package agent

import (
	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/ai"
	"github.com/clivern/actx0/pkg/qdrant"
)

// Dependencies are the collaborators required by the agent service.
type Dependencies struct {
	Embed         *ai.EmbedClient
	Vectors       *qdrant.Client
	LiteLLM       *ai.ChatClient
	LargeLLM      *ai.ChatClient
	Agents        db.AgentRepository
	Sessions      db.AgentSessionRepository
	Messages      db.SessionMessageRepository
	Memories      db.SessionMemoryRepository
	Usage         db.UsageRepository
	Subscriptions db.SubscriptionRepository
}

// Service indexes session content and runs semantic search over it.
type Service struct {
	embed         *ai.EmbedClient
	vectors       *qdrant.Client
	liteLLM       *ai.ChatClient
	largeLLM      *ai.ChatClient
	agents        db.AgentRepository
	sessions      db.AgentSessionRepository
	messages      db.SessionMessageRepository
	memories      db.SessionMemoryRepository
	usage         db.UsageRepository
	subscriptions db.SubscriptionRepository
}

// New returns an agent service.
func New(deps Dependencies) *Service {
	return &Service{
		embed:         deps.Embed,
		vectors:       deps.Vectors,
		liteLLM:       deps.LiteLLM,
		largeLLM:      deps.LargeLLM,
		agents:        deps.Agents,
		sessions:      deps.Sessions,
		messages:      deps.Messages,
		memories:      deps.Memories,
		usage:         deps.Usage,
		subscriptions: deps.Subscriptions,
	}
}
