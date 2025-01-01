// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/clivern/actx0/conf"
	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/migration"
	"github.com/clivern/actx0/pkg/ai"
	"github.com/clivern/actx0/pkg/qdrant"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// ConsolidateMemory summarizes and compacts session memory.
func (s *Service) ConsolidateMemory(ctx context.Context, req ConsolidateMemoryRequest) error {
	return ErrNotImplemented
}

// StoreMemory saves a memory item for a session.
func (s *Service) StoreMemory(ctx context.Context, req StoreMemoryRequest) (*db.SessionMemory, error) {
	return nil, ErrNotImplemented
}

// RecallMemory returns memories relevant to the query.
func (s *Service) RecallMemory(ctx context.Context, req RecallMemoryRequest) ([]MemoryHit, error) {
	return nil, ErrNotImplemented
}

// ListMemories returns stored memories for a session.
func (s *Service) ListMemories(ctx context.Context, sessionId db.Id, limit, offset int) ([]*db.SessionMemory, error) {
	return nil, ErrNotImplemented
}

// IndexSessionMessage embeds a session message into Qdrant.
func (s *Service) IndexSessionMessage(ctx context.Context, messageId db.Id) error {
	message, err := s.messages.GetById(messageId)
	if err != nil {
		return fmt.Errorf("%w: get message: %v", ErrIndexFailed, err)
	}
	if message == nil {
		return ErrMessageNotFound
	}

	session, err := s.sessions.GetById(message.SessionId)
	if err != nil {
		return fmt.Errorf("%w: get session: %v", ErrIndexFailed, err)
	}
	if session == nil {
		return ErrSessionNotFound
	}

	embeddings, usage, err := s.embed.EmbedDocuments(ctx, []string{message.Content}, ai.EmbedOptions{
		Dimensions: int64(viper.GetInt("app.ai.embed.qdb.vector_size")),
	})
	if err != nil {
		return fmt.Errorf("%w: embed message: %v", ErrIndexFailed, err)
	}

	sub, err := s.subscriptions.GetByWorkspaceId(session.WorkspaceId)
	if err != nil {
		log.Warn().
			Err(err).
			Str("workspace_id", session.WorkspaceId.String()).
			Int64("total_tokens", usage.TotalTokens).
			Int64("total_cost", usage.Cost).
			Msg("Failed to record AI usage")
	}

	if sub != nil && sub.CurrentPeriodStart != nil && sub.CurrentPeriodEnd != nil {
		s.usage.IncrementByPeriod(
			session.WorkspaceId,
			db.UsageTypeAITokens,
			sub.CurrentPeriodStart.UTC(),
			sub.CurrentPeriodEnd.UTC(),
			usage.TotalTokens,
			db.UsageUnitTokens,
		)

		s.usage.IncrementByPeriod(
			session.WorkspaceId,
			db.UsageTypeAICost,
			sub.CurrentPeriodStart.UTC(),
			sub.CurrentPeriodEnd.UTC(),
			usage.Cost,
			db.UsageUnitNanoUSD,
		)
	}

	log.Info().
		Str("message_id", messageId.String()).
		Int64("prompt_tokens", usage.PromptTokens).
		Int64("total_tokens", usage.TotalTokens).
		Int64("total_cost", usage.Cost).
		Msg("Session message embedded")

	err = s.vectors.DeleteByFilter(ctx, migration.AgentSessionMessagesCollection, map[string]string{
		migration.PayloadMessageID: message.Id.String(),
	})
	if err != nil {
		return fmt.Errorf("%w: delete existing vectors: %v", ErrIndexFailed, err)
	}

	err = s.vectors.Upsert(ctx, migration.AgentSessionMessagesCollection, []qdrant.Point{
		{
			Id:     message.InternalId.String(),
			Vector: qdrant.Float32Vector(embeddings[0]),
			Payload: map[string]string{
				migration.PayloadWorkspaceID: session.WorkspaceId.String(),
				migration.PayloadAgentID:     session.AgentId.String(),
				migration.PayloadSessionID:   session.Id.String(),
				migration.PayloadMessageID:   message.Id.String(),
				migration.PayloadText:        message.Content,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("%w: upsert vectors: %v", ErrIndexFailed, err)
	}

	return nil
}

// IndexSessionMemory embeds a session memory into Qdrant.
func (s *Service) IndexSessionMemory(ctx context.Context, memoryId db.Id) error {
	memory, err := s.memories.GetById(memoryId)
	if err != nil {
		return fmt.Errorf("%w: get memory: %v", ErrIndexFailed, err)
	}
	if memory == nil {
		return ErrMemoryNotFound
	}

	session, err := s.sessions.GetById(memory.SessionId)
	if err != nil {
		return fmt.Errorf("%w: get session: %v", ErrIndexFailed, err)
	}
	if session == nil {
		return ErrSessionNotFound
	}

	embeddings, usage, err := s.embed.EmbedDocuments(ctx, []string{memory.Content}, ai.EmbedOptions{
		Dimensions: int64(viper.GetInt("app.ai.embed.qdb.vector_size")),
	})
	if err != nil {
		return fmt.Errorf("%w: embed memory: %v", ErrIndexFailed, err)
	}

	sub, err := s.subscriptions.GetByWorkspaceId(session.WorkspaceId)
	if err != nil {
		log.Warn().
			Err(err).
			Str("workspace_id", session.WorkspaceId.String()).
			Int64("total_tokens", usage.TotalTokens).
			Int64("total_cost", usage.Cost).
			Msg("Failed to record AI usage")
	}

	if sub != nil && sub.CurrentPeriodStart != nil && sub.CurrentPeriodEnd != nil {
		s.usage.IncrementByPeriod(
			session.WorkspaceId,
			db.UsageTypeAITokens,
			sub.CurrentPeriodStart.UTC(),
			sub.CurrentPeriodEnd.UTC(),
			usage.TotalTokens,
			db.UsageUnitTokens,
		)

		s.usage.IncrementByPeriod(
			session.WorkspaceId,
			db.UsageTypeAICost,
			sub.CurrentPeriodStart.UTC(),
			sub.CurrentPeriodEnd.UTC(),
			usage.Cost,
			db.UsageUnitNanoUSD,
		)
	}

	log.Info().
		Str("memory_id", memoryId.String()).
		Int64("prompt_tokens", usage.PromptTokens).
		Int64("total_tokens", usage.TotalTokens).
		Int64("total_cost", usage.Cost).
		Msg("Session memory embedded")

	err = s.vectors.DeleteByFilter(ctx, migration.AgentSessionMemoriesCollection, map[string]string{
		migration.PayloadMemoryID: memory.Id.String(),
	})
	if err != nil {
		return fmt.Errorf("%w: delete existing vectors: %v", ErrIndexFailed, err)
	}

	err = s.vectors.Upsert(ctx, migration.AgentSessionMemoriesCollection, []qdrant.Point{
		{
			Id:     memory.InternalId.String(),
			Vector: qdrant.Float32Vector(embeddings[0]),
			Payload: map[string]string{
				migration.PayloadWorkspaceID: session.WorkspaceId.String(),
				migration.PayloadAgentID:     session.AgentId.String(),
				migration.PayloadSessionID:   session.Id.String(),
				migration.PayloadMemoryID:    memory.Id.String(),
				migration.PayloadText:        memory.Content,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("%w: upsert vectors: %v", ErrIndexFailed, err)
	}

	return nil
}

// GenerateMemoryFromMessage extracts and stores memory from a session message.
func (s *Service) GenerateMemoryFromMessage(ctx context.Context, messageId db.Id) error {
	message, err := s.messages.GetById(messageId)
	if err != nil {
		return fmt.Errorf("%w: get message: %v", ErrGenerateFailed, err)
	}
	if message == nil {
		return ErrMessageNotFound
	}

	session, err := s.sessions.GetById(message.SessionId)
	if err != nil {
		return fmt.Errorf("%w: get session: %v", ErrGenerateFailed, err)
	}
	if session == nil {
		return ErrSessionNotFound
	}

	reply, usage, err := s.liteLLM.Complete(ctx, []ai.Message{
		{Role: "system", Content: conf.GenerateMemorySystemPrompt},
		{Role: "user", Content: fmt.Sprintf("role=%s\n%s", message.Role, message.Content)},
	})
	if err != nil {
		return fmt.Errorf("%w: complete: %v", ErrGenerateFailed, err)
	}

	sub, err := s.subscriptions.GetByWorkspaceId(session.WorkspaceId)
	if err != nil {
		log.Warn().
			Err(err).
			Str("workspace_id", session.WorkspaceId.String()).
			Int64("total_tokens", usage.TotalTokens).
			Int64("total_cost", usage.Cost).
			Msg("Failed to record AI usage")
	}

	if sub != nil && sub.CurrentPeriodStart != nil && sub.CurrentPeriodEnd != nil {
		s.usage.IncrementByPeriod(
			session.WorkspaceId,
			db.UsageTypeAITokens,
			sub.CurrentPeriodStart.UTC(),
			sub.CurrentPeriodEnd.UTC(),
			usage.TotalTokens,
			db.UsageUnitTokens,
		)

		s.usage.IncrementByPeriod(
			session.WorkspaceId,
			db.UsageTypeAICost,
			sub.CurrentPeriodStart.UTC(),
			sub.CurrentPeriodEnd.UTC(),
			usage.Cost,
			db.UsageUnitNanoUSD,
		)
	}

	var result struct {
		Skip    bool   `json:"skip"`
		Kind    string `json:"kind"`
		Content string `json:"content"`
	}
	err = json.Unmarshal([]byte(strings.TrimSpace(reply)), &result)
	if err != nil {
		return fmt.Errorf("%w: parse reply: %v", ErrGenerateFailed, err)
	}

	result.Kind = strings.TrimSpace(result.Kind)
	result.Content = strings.TrimSpace(result.Content)

	if result.Skip || result.Content == "" {
		return nil
	}
	if result.Kind != MemoryKindShortLived && result.Kind != MemoryKindLongLived {
		return nil
	}

	memory := &db.SessionMemory{
		SessionId: session.Id,
		Kind:      result.Kind,
		Content:   result.Content,
	}
	err = s.memories.Create(memory)
	if err != nil {
		return fmt.Errorf("%w: create memory: %v", ErrGenerateFailed, err)
	}

	err = s.IndexSessionMemory(ctx, memory.Id)
	if err != nil {
		return fmt.Errorf("%w: index memory: %v", ErrGenerateFailed, err)
	}

	return nil
}

// DeleteSessionMemory deletes vectors for one memory item.
func (s *Service) DeleteSessionMemory(ctx context.Context, memoryId db.Id) error {
	err := s.vectors.DeleteByFilter(ctx, migration.AgentSessionMemoriesCollection, map[string]string{
		migration.PayloadMemoryID: memoryId.String(),
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("memory_id", memoryId.String()).
			Msg("Failed to delete memory vectors")
		return fmt.Errorf("%w: delete memory vectors: %v", ErrDeleteFailed, err)
	}

	return nil
}

// DeleteSessionMemories deletes all memory vectors for a session.
func (s *Service) DeleteSessionMemories(ctx context.Context, agentId, sessionId db.Id) error {
	err := s.vectors.DeleteByFilter(ctx, migration.AgentSessionMemoriesCollection, map[string]string{
		migration.PayloadAgentID:   agentId.String(),
		migration.PayloadSessionID: sessionId.String(),
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("agent_id", agentId.String()).
			Str("session_id", sessionId.String()).
			Msg("Failed to delete session memory vectors")
		return fmt.Errorf("%w: delete memory vectors: %v", ErrDeleteFailed, err)
	}

	return nil
}

// DeleteSessionMessages deletes all message vectors for a session.
func (s *Service) DeleteSessionMessages(ctx context.Context, agentId, sessionId db.Id) error {
	err := s.vectors.DeleteByFilter(ctx, migration.AgentSessionMessagesCollection, map[string]string{
		migration.PayloadAgentID:   agentId.String(),
		migration.PayloadSessionID: sessionId.String(),
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("agent_id", agentId.String()).
			Str("session_id", sessionId.String()).
			Msg("Failed to delete session message vectors")
		return fmt.Errorf("%w: delete message vectors: %v", ErrDeleteFailed, err)
	}

	return nil
}

// DeleteSessionMessage deletes vectors for one message.
func (s *Service) DeleteSessionMessage(ctx context.Context, messageId db.Id) error {
	err := s.vectors.DeleteByFilter(ctx, migration.AgentSessionMessagesCollection, map[string]string{
		migration.PayloadMessageID: messageId.String(),
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("message_id", messageId.String()).
			Msg("Failed to delete message vectors")
		return fmt.Errorf("%w: delete message vectors: %v", ErrDeleteFailed, err)
	}

	return nil
}

// DeleteAgentSession deletes all message and memory vectors for a session.
func (s *Service) DeleteAgentSession(ctx context.Context, agentId, sessionId db.Id) error {
	err := s.DeleteSessionMessages(ctx, agentId, sessionId)
	if err != nil {
		return err
	}

	return s.DeleteSessionMemories(ctx, agentId, sessionId)
}

// DeleteAgent removes all session message and memory vectors for an agent.
func (s *Service) DeleteAgent(ctx context.Context, agentId db.Id) error {
	err := s.vectors.DeleteByFilter(ctx, migration.AgentSessionMessagesCollection, map[string]string{
		migration.PayloadAgentID: agentId.String(),
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("agent_id", agentId.String()).
			Msg("Failed to delete agent message vectors")
		return fmt.Errorf("%w: delete message vectors: %v", ErrDeleteFailed, err)
	}

	err = s.vectors.DeleteByFilter(ctx, migration.AgentSessionMemoriesCollection, map[string]string{
		migration.PayloadAgentID: agentId.String(),
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("agent_id", agentId.String()).
			Msg("Failed to delete agent memory vectors")
		return fmt.Errorf("%w: delete memory vectors: %v", ErrDeleteFailed, err)
	}

	return nil
}
