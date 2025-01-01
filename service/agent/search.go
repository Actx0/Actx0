// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package agent

import (
	"context"
	"fmt"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/migration"
	"github.com/clivern/actx0/pkg/ai"
	"github.com/clivern/actx0/pkg/qdrant"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// SearchSessionMessages returns session messages relevant to the query.
func (s *Service) SearchSessionMessages(ctx context.Context, req SearchSessionRequest) ([]MessageHit, ai.Usage, error) {
	vector, usage, err := s.embed.EmbedQuery(ctx, req.Query, ai.EmbedOptions{
		Dimensions: int64(viper.GetInt("app.ai.embed.qdb.vector_size")),
	})
	if err != nil {
		return nil, usage, fmt.Errorf("%w: embed query: %v", ErrSearchFailed, err)
	}

	results, err := s.vectors.Search(ctx, migration.AgentSessionMessagesCollection, qdrant.Query{
		Vector: qdrant.Float32Vector(vector),
		Filters: map[string]string{
			migration.PayloadWorkspaceID: req.WorkspaceId.String(),
			migration.PayloadAgentID:     req.AgentId.String(),
			migration.PayloadSessionID:   req.SessionId.String(),
		},
		Limit: req.Limit,
	})
	if err != nil {
		return nil, usage, fmt.Errorf("%w: vector search: %v", ErrSearchFailed, err)
	}

	hits := make([]MessageHit, 0, len(results))
	for _, result := range results {
		message, err := s.messages.GetById(db.Id(result.Payload[migration.PayloadMessageID]))
		if err != nil || message == nil {
			continue
		}

		hits = append(hits, MessageHit{
			MessageId: message.Id,
			Role:      message.Role,
			Score:     result.Score,
			Text:      result.Payload[migration.PayloadText],
		})
	}

	log.Info().
		Str("session_id", req.SessionId.String()).
		Int("hits", len(hits)).
		Int64("prompt_tokens", usage.PromptTokens).
		Int64("total_tokens", usage.TotalTokens).
		Int64("total_cost", usage.Cost).
		Msg("Session message search completed")

	return hits, usage, nil
}

// SearchSessionMemories returns session memories relevant to the query.
func (s *Service) SearchSessionMemories(ctx context.Context, req SearchSessionRequest) ([]MemoryHit, ai.Usage, error) {
	vector, usage, err := s.embed.EmbedQuery(ctx, req.Query, ai.EmbedOptions{
		Dimensions: int64(viper.GetInt("app.ai.embed.qdb.vector_size")),
	})
	if err != nil {
		return nil, usage, fmt.Errorf("%w: embed query: %v", ErrSearchFailed, err)
	}

	results, err := s.vectors.Search(ctx, migration.AgentSessionMemoriesCollection, qdrant.Query{
		Vector: qdrant.Float32Vector(vector),
		Filters: map[string]string{
			migration.PayloadWorkspaceID: req.WorkspaceId.String(),
			migration.PayloadAgentID:     req.AgentId.String(),
			migration.PayloadSessionID:   req.SessionId.String(),
		},
		Limit: req.Limit,
	})
	if err != nil {
		return nil, usage, fmt.Errorf("%w: vector search: %v", ErrSearchFailed, err)
	}

	hits := make([]MemoryHit, 0, len(results))
	for _, result := range results {
		memory, err := s.memories.GetById(db.Id(result.Payload[migration.PayloadMemoryID]))
		if err != nil || memory == nil {
			continue
		}

		hits = append(hits, MemoryHit{
			MemoryId: memory.Id,
			Kind:     memory.Kind,
			Content:  result.Payload[migration.PayloadText],
			Score:    result.Score,
		})
	}

	log.Info().
		Str("session_id", req.SessionId.String()).
		Int("hits", len(hits)).
		Int64("prompt_tokens", usage.PromptTokens).
		Int64("total_tokens", usage.TotalTokens).
		Int64("total_cost", usage.Cost).
		Msg("Session memory search completed")

	return hits, usage, nil
}
