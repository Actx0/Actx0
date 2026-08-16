// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/stripe"
	"github.com/clivern/actx0/pkg/util"
	"github.com/clivern/actx0/service/agent"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

var (
	ErrSessionMemoryNotFound       = errors.New("session memory not found")
	ErrFailedCreateSessionMemory   = errors.New("failed create session memory")
	ErrFailedListSessionMemories   = errors.New("failed list session memories")
	ErrFailedGetSessionMemory      = errors.New("failed get session memory")
	ErrFailedUpdateSessionMemory   = errors.New("failed update session memory")
	ErrFailedDeleteSessionMemory   = errors.New("failed delete session memory")
	ErrFailedSearchSessionMemories = errors.New("failed search session memories")

	ErrFailedRegenerateSessionMemories = errors.New("failed regenerate session memories")
	ErrMemoryPipelineDisabled          = errors.New("memory pipeline disabled for agent")
)

// CreateSessionMemoryRequest is the body for creating a session memory.
type CreateSessionMemoryRequest struct {
	Kind    string  `json:"kind" validate:"required,oneof=summary fact preference short_lived long_lived" label:"Kind"`
	Content string  `json:"content" validate:"required" label:"Content"`
	Meta    *string `json:"meta" validate:"omitempty,json" label:"Meta"`
}

// BatchCreateSessionMemoriesRequest is the body for batch creating session memories.
type BatchCreateSessionMemoriesRequest struct {
	Memories []CreateSessionMemoryRequest `json:"memories" validate:"required,min=1,max=100,dive" label:"Memories"`
}

// BatchDeleteSessionMemoriesRequest is the body for batch deleting session memories.
type BatchDeleteSessionMemoriesRequest struct {
	IDs []string `json:"ids" validate:"required,min=1,max=100,dive,required" label:"IDs"`
}

// UpdateSessionMemoryRequest is the body for updating a session memory.
type UpdateSessionMemoryRequest struct {
	Kind    string  `json:"kind" validate:"omitempty,oneof=summary fact preference short_lived long_lived" label:"Kind"`
	Content string  `json:"content" validate:"required" label:"Content"`
	Meta    *string `json:"meta" validate:"omitempty,json" label:"Meta"`
}

// SearchSessionMemoriesRequest is the body for searching session memories.
type SearchSessionMemoriesRequest struct {
	Query string `json:"query" validate:"required,min=1" label:"Query"`
	Limit int    `json:"limit" validate:"omitempty,min=1,max=100" label:"Limit"`
}

// SearchSessionMemoryHit is a scored session memory for API responses.
type SearchSessionMemoryHit struct {
	Id    db.Id   `json:"id"`
	Kind  string  `json:"kind"`
	Score float32 `json:"score"`
	Text  string  `json:"text"`
}

// SearchSessionMemoriesResponse is returned when searching session memories.
type SearchSessionMemoriesResponse struct {
	Results []SearchSessionMemoryHit `json:"results"`
}

// SessionMemoryResponse is a session memory shaped for API responses.
type SessionMemoryResponse struct {
	Id        db.Id           `json:"id"`
	SessionId db.Id           `json:"sessionId"`
	Kind      string          `json:"kind"`
	Content   string          `json:"content"`
	Meta      json.RawMessage `json:"meta,omitempty"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

// ListSessionMemoriesResponse is returned when listing session memories.
type ListSessionMemoriesResponse struct {
	Memories []*SessionMemoryResponse
	Total    int64
}

// RegenerateSessionMemoriesResponse is returned when regenerating session memories.
type RegenerateSessionMemoriesResponse struct {
	Queued int `json:"queued"`
}

// Memory is the module for session memory CRUD.
type Memory struct {
	MemoryRepository       db.SessionMemoryRepository
	MessageRepository      db.SessionMessageRepository
	AgentSession           *AgentSession
	Agent                  *agent.Service
	UsageRepository        db.UsageRepository
	SubscriptionRepository db.SubscriptionRepository
}

// NewMemory creates a session memory module with the given repositories.
func NewMemory(
	memories db.SessionMemoryRepository,
	messages db.SessionMessageRepository,
	agentSessions *AgentSession,
	subscriptions db.SubscriptionRepository,
) *Memory {
	return &Memory{
		MemoryRepository:       memories,
		MessageRepository:      messages,
		AgentSession:           agentSessions,
		SubscriptionRepository: subscriptions,
	}
}

// CreateSessionMemory creates a memory in a session.
func (m *Memory) CreateSessionMemory(workspaceId, agentId, sessionId db.Id, req *CreateSessionMemoryRequest) (*SessionMemoryResponse, error) {
	session, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	reached, err := NewUsage().IsMemoryRecordsReached(
		m.MemoryRepository,
		m.SubscriptionRepository,
		workspaceId,
		1,
	)
	if err != nil {
		return nil, err
	}
	if reached {
		return nil, ErrWorkspaceCapacityReached
	}

	memory := &db.SessionMemory{
		SessionId: session.Id,
		Kind:      req.Kind,
		Content:   req.Content,
		Meta:      req.Meta,
	}

	err = m.MemoryRepository.Create(memory)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateSessionMemory, err)
	}

	subscription, err := m.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	if subscription != nil && subscription.Plan != stripe.PlanHobby {
		err = m.EnqueueIndexTask(workspaceId, memory.Id)
		if err != nil {
			return nil, err
		}
	}

	return &SessionMemoryResponse{
		Id:        memory.Id,
		SessionId: memory.SessionId,
		Kind:      memory.Kind,
		Content:   memory.Content,
		Meta:      util.JSONRawFromString(memory.Meta),
		CreatedAt: memory.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: memory.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// BatchCreateSessionMemories creates multiple memories in a session.
func (m *Memory) BatchCreateSessionMemories(workspaceId, agentId, sessionId db.Id, req *BatchCreateSessionMemoriesRequest) ([]*SessionMemoryResponse, error) {
	session, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	reached, err := NewUsage().IsMemoryRecordsReached(
		m.MemoryRepository,
		m.SubscriptionRepository,
		workspaceId,
		int64(len(req.Memories)),
	)
	if err != nil {
		return nil, err
	}
	if reached {
		return nil, ErrWorkspaceCapacityReached
	}

	subscription, err := m.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	indexEnabled := subscription != nil && subscription.Plan != stripe.PlanHobby

	responses := make([]*SessionMemoryResponse, 0, len(req.Memories))
	for _, item := range req.Memories {
		memory := &db.SessionMemory{
			SessionId: session.Id,
			Kind:      item.Kind,
			Content:   item.Content,
			Meta:      item.Meta,
		}

		err = m.MemoryRepository.Create(memory)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedCreateSessionMemory, err)
		}

		if indexEnabled {
			err = m.EnqueueIndexTask(workspaceId, memory.Id)
			if err != nil {
				return nil, err
			}
		}

		responses = append(responses, &SessionMemoryResponse{
			Id:        memory.Id,
			SessionId: memory.SessionId,
			Kind:      memory.Kind,
			Content:   memory.Content,
			Meta:      util.JSONRawFromString(memory.Meta),
			CreatedAt: memory.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: memory.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return responses, nil
}

// GetSessionMemory returns one memory by id.
func (m *Memory) GetSessionMemory(workspaceId, agentId, sessionId, memoryId db.Id) (*SessionMemoryResponse, error) {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	memory, err := m.MemoryRepository.GetById(memoryId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetSessionMemory, err)
	}
	if memory == nil || memory.SessionId != sessionId {
		return nil, ErrSessionMemoryNotFound
	}

	return &SessionMemoryResponse{
		Id:        memory.Id,
		SessionId: memory.SessionId,
		Kind:      memory.Kind,
		Content:   memory.Content,
		Meta:      util.JSONRawFromString(memory.Meta),
		CreatedAt: memory.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: memory.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// UpdateSessionMemory updates a memory by id.
func (m *Memory) UpdateSessionMemory(workspaceId, agentId, sessionId, memoryId db.Id, req *UpdateSessionMemoryRequest) (*SessionMemoryResponse, error) {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	memory, err := m.MemoryRepository.GetById(memoryId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdateSessionMemory, err)
	}
	if memory == nil || memory.SessionId != sessionId {
		return nil, ErrSessionMemoryNotFound
	}

	if lo.IsNotEmpty(req.Kind) {
		memory.Kind = req.Kind
	}
	memory.Content = req.Content
	memory.Meta = req.Meta

	err = m.MemoryRepository.Update(memory)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdateSessionMemory, err)
	}

	subscription, err := m.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	if subscription != nil && subscription.Plan != stripe.PlanHobby {
		err = m.EnqueueIndexTask(workspaceId, memory.Id)
		if err != nil {
			return nil, err
		}
	}

	return &SessionMemoryResponse{
		Id:        memory.Id,
		SessionId: memory.SessionId,
		Kind:      memory.Kind,
		Content:   memory.Content,
		Meta:      util.JSONRawFromString(memory.Meta),
		CreatedAt: memory.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: memory.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// DeleteSessionMemory deletes a memory by id.
func (m *Memory) DeleteSessionMemory(workspaceId, agentId, sessionId, memoryId db.Id) error {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return err
	}

	memory, err := m.MemoryRepository.GetById(memoryId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteSessionMemory, err)
	}
	if memory == nil || memory.SessionId != sessionId {
		return ErrSessionMemoryNotFound
	}

	err = m.MemoryRepository.Delete(memoryId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteSessionMemory, err)
	}

	return m.EnqueueDeleteTask(workspaceId, memoryId)
}

// BatchDeleteSessionMemories deletes multiple memories by id.
func (m *Memory) BatchDeleteSessionMemories(workspaceId, agentId, sessionId db.Id, req *BatchDeleteSessionMemoriesRequest) error {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return err
	}

	for _, id := range req.IDs {
		memoryId := db.Id(id)

		memory, err := m.MemoryRepository.GetById(memoryId)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFailedDeleteSessionMemory, err)
		}
		if memory == nil || memory.SessionId != sessionId {
			return ErrSessionMemoryNotFound
		}

		err = m.MemoryRepository.Delete(memoryId)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFailedDeleteSessionMemory, err)
		}

		err = m.EnqueueDeleteTask(workspaceId, memoryId)
		if err != nil {
			return err
		}
	}

	return nil
}

// ListSessionMemories returns paginated memories for a session.
func (m *Memory) ListSessionMemories(workspaceId, agentId, sessionId db.Id, limit, offset int) (*ListSessionMemoriesResponse, error) {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	total, err := m.MemoryRepository.CountBySessionId(sessionId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListSessionMemories, err)
	}

	memories, err := m.MemoryRepository.ListBySessionId(sessionId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListSessionMemories, err)
	}

	list := make([]*SessionMemoryResponse, 0, len(memories))
	for _, item := range memories {
		list = append(list, &SessionMemoryResponse{
			Id:        item.Id,
			SessionId: item.SessionId,
			Kind:      item.Kind,
			Content:   item.Content,
			Meta:      util.JSONRawFromString(item.Meta),
			CreatedAt: item.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: item.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return &ListSessionMemoriesResponse{Memories: list, Total: total}, nil
}

// SearchSessionMemories finds relevant session memories for a semantic query.
func (m *Memory) SearchSessionMemories(ctx context.Context, workspaceId, agentId, sessionId db.Id, query string, limit int) (*SearchSessionMemoriesResponse, error) {
	session, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	subscription, err := m.SubscriptionRepository.GetByWorkspaceId(workspaceId)
	if err != nil {
		return nil, err
	}
	if subscription == nil || subscription.Plan == stripe.PlanHobby {
		return nil, ErrSemanticSearchRequiresPaidPlan
	}

	usageMod := NewUsage()

	reached, err := usageMod.IsAITokensReached(m.UsageRepository, m.SubscriptionRepository, workspaceId, 1)
	if err != nil {
		return nil, err
	}
	if reached {
		return nil, ErrWorkspaceCapacityReached
	}

	reached, err = usageMod.IsAICostReached(m.UsageRepository, m.SubscriptionRepository, workspaceId, 0)
	if err != nil {
		return nil, err
	}
	if reached {
		return nil, ErrWorkspaceCapacityReached
	}

	hits, usage, err := m.Agent.SearchSessionMemories(ctx, agent.SearchSessionRequest{
		WorkspaceId: workspaceId,
		AgentId:     agentId,
		SessionId:   session.Id,
		Query:       query,
		Limit:       uint64(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedSearchSessionMemories, err)
	}

	err = usageMod.IncrementAIUsage(
		m.UsageRepository,
		m.SubscriptionRepository,
		workspaceId,
		usage.TotalTokens,
		usage.Cost,
	)
	if err != nil {
		if errors.Is(err, ErrWorkspaceCapacityReached) {
			return nil, err
		}

		util.CaptureError(err)
		log.Warn().
			Err(err).
			Str("workspaceId", workspaceId.String()).
			Int64("total_tokens", usage.TotalTokens).
			Int64("total_cost", usage.Cost).
			Msg("Failed to record AI usage for session memory search")
	}

	results := make([]SearchSessionMemoryHit, 0, len(hits))
	for _, hit := range hits {
		results = append(results, SearchSessionMemoryHit{
			Id:    hit.MemoryId,
			Kind:  hit.Kind,
			Score: hit.Score,
			Text:  hit.Content,
		})
	}

	return &SearchSessionMemoriesResponse{
		Results: results,
	}, nil
}

// RegenerateSessionMemories queues a task that clears and rebuilds session memories.
func (m *Memory) RegenerateSessionMemories(workspaceId, agentId, sessionId db.Id) (*RegenerateSessionMemoriesResponse, error) {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	agent, err := m.AgentSession.AgentRepository.GetById(agentId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedRegenerateSessionMemories, err)
	}
	if agent == nil || agent.WorkspaceId != workspaceId {
		return nil, ErrAgentNotFound
	}

	var cfg AgentConfigs
	_ = json.Unmarshal([]byte(*agent.Config), &cfg)
	if !cfg.MemoryPipeline {
		return nil, ErrMemoryPipelineDisabled
	}

	err = m.MemoryRepository.DeleteBySessionId(sessionId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedRegenerateSessionMemories, err)
	}

	count, err := m.MessageRepository.CountBySessionId(sessionId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedRegenerateSessionMemories, err)
	}

	err = m.EnqueueRegenerateMemoriesTask(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	return &RegenerateSessionMemoriesResponse{Queued: int(count)}, nil
}

// EnqueueIndexTask enqueues an async task to index a session memory.
func (m *Memory) EnqueueIndexTask(workspaceId, memoryId db.Id) error {
	_, err := GetAsyncMgr().CreateTask(&CreateAsyncTaskOptions{
		WorkspaceId: workspaceId,
		Type:        "task.agent.session.memory.index",
		Priority:    AsyncPriorityLow,
		Payload: map[string]string{
			"memoryId": memoryId.String(),
		},
	})

	return err
}

// EnqueueDeleteTask enqueues an async task to delete session memory vectors.
func (m *Memory) EnqueueDeleteTask(workspaceId, memoryId db.Id) error {
	_, err := GetAsyncMgr().CreateTask(&CreateAsyncTaskOptions{
		WorkspaceId: workspaceId,
		Type:        "task.agent.session.memory.delete",
		Priority:    AsyncPriorityLow,
		Payload: map[string]string{
			"memoryId": memoryId.String(),
		},
	})

	return err
}

// EnqueueRegenerateMemoriesTask enqueues an async task to regenerate session memories.
func (m *Memory) EnqueueRegenerateMemoriesTask(workspaceId, agentId, sessionId db.Id) error {
	_, err := GetAsyncMgr().CreateTask(&CreateAsyncTaskOptions{
		WorkspaceId: workspaceId,
		Type:        "task.agent.session.memories.regenerate",
		Priority:    AsyncPriorityLow,
		Payload: map[string]string{
			"agentId":   agentId.String(),
			"sessionId": sessionId.String(),
		},
	})

	return err
}
