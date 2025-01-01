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
	"github.com/clivern/actx0/pkg/util"
	"github.com/clivern/actx0/service/agent"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

var (
	ErrSessionMessageNotFound      = errors.New("session message not found")
	ErrFailedCreateSessionMessage  = errors.New("failed create session message")
	ErrFailedListSessionMessages   = errors.New("failed list session messages")
	ErrFailedGetSessionMessage     = errors.New("failed get session message")
	ErrFailedUpdateSessionMessage  = errors.New("failed update session message")
	ErrFailedDeleteSessionMessage  = errors.New("failed delete session message")
	ErrFailedSearchSessionMessages = errors.New("failed search session messages")
)

// CreateSessionMessageRequest is the body for creating a session message.
type CreateSessionMessageRequest struct {
	Role    string  `json:"role" validate:"required,oneof=system user assistant" label:"Role"`
	Content string  `json:"content" validate:"required" label:"Content"`
	Meta    *string `json:"meta" validate:"omitempty,json" label:"Meta"`
}

// BatchCreateSessionMessagesRequest is the body for batch creating session messages.
type BatchCreateSessionMessagesRequest struct {
	Messages []CreateSessionMessageRequest `json:"messages" validate:"required,min=1,max=100,dive" label:"Messages"`
}

// BatchDeleteSessionMessagesRequest is the body for batch deleting session messages.
type BatchDeleteSessionMessagesRequest struct {
	IDs []string `json:"ids" validate:"required,min=1,max=100,dive,required" label:"IDs"`
}

// UpdateSessionMessageRequest is the body for updating a session message.
type UpdateSessionMessageRequest struct {
	Role    string  `json:"role" validate:"omitempty,oneof=system user assistant" label:"Role"`
	Content string  `json:"content" validate:"required" label:"Content"`
	Meta    *string `json:"meta" validate:"omitempty,json" label:"Meta"`
}

// SearchSessionMessagesRequest is the body for searching session messages.
type SearchSessionMessagesRequest struct {
	Query string `json:"query" validate:"required,min=1" label:"Query"`
	Limit int    `json:"limit" validate:"omitempty,min=1,max=100" label:"Limit"`
}

// SearchSessionMessageHit is a scored session message for API responses.
type SearchSessionMessageHit struct {
	Id    db.Id   `json:"id"`
	Role  string  `json:"role"`
	Score float32 `json:"score"`
	Text  string  `json:"text"`
}

// SearchSessionMessagesResponse is returned when searching session messages.
type SearchSessionMessagesResponse struct {
	Results []SearchSessionMessageHit `json:"results"`
}

// SessionMessageResponse is a session message shaped for API responses.
type SessionMessageResponse struct {
	Id        db.Id           `json:"id"`
	SessionId db.Id           `json:"sessionId"`
	Role      string          `json:"role"`
	Content   string          `json:"content"`
	Meta      json.RawMessage `json:"meta,omitempty"`
	CreatedAt string          `json:"createdAt"`
}

// ListSessionMessagesResponse is returned when listing session messages.
type ListSessionMessagesResponse struct {
	Messages []*SessionMessageResponse
	Total    int64
}

// Message is the module for session message CRUD.
type Message struct {
	MessageRepository      db.SessionMessageRepository
	SessionRepository      db.AgentSessionRepository
	Agent                  *agent.Service
	UsageRepository        db.UsageRepository
	SubscriptionRepository db.SubscriptionRepository
}

// NewMessage creates a session message module with the given repositories.
func NewMessage(messages db.SessionMessageRepository, sessions db.AgentSessionRepository) *Message {
	return &Message{
		MessageRepository: messages,
		SessionRepository: sessions,
	}
}

func (m *Message) loadAgentSession(workspaceId, agentId, sessionId db.Id) (*db.AgentSession, error) {
	session, err := m.SessionRepository.GetById(sessionId)
	if err != nil {
		return nil, err
	}
	if session == nil || session.AgentId != agentId || session.WorkspaceId != workspaceId {
		return nil, ErrAgentSessionNotFound
	}

	return session, nil
}

// CreateSessionMessage creates a message in a session.
func (m *Message) CreateSessionMessage(workspaceId, agentId, sessionId db.Id, req *CreateSessionMessageRequest) (*SessionMessageResponse, error) {
	session, err := m.loadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	message := &db.SessionMessage{
		SessionId: session.Id,
		Role:      req.Role,
		Content:   req.Content,
		Meta:      req.Meta,
	}

	err = m.MessageRepository.Create(message)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateSessionMessage, err)
	}

	err = m.SessionRepository.TouchLastActivityAt(session.Id, message.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateSessionMessage, err)
	}

	err = m.EnqueueIndexTask(workspaceId, message.Id)
	if err != nil {
		return nil, err
	}

	if message.Role == "user" {
		err = m.EnqueueGenerateMemoryTask(workspaceId, message.Id)
		if err != nil {
			return nil, err
		}
	}

	return &SessionMessageResponse{
		Id:        message.Id,
		SessionId: message.SessionId,
		Role:      message.Role,
		Content:   message.Content,
		Meta:      util.JSONRawFromString(message.Meta),
		CreatedAt: message.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// BatchCreateSessionMessages creates multiple messages in a session.
func (m *Message) BatchCreateSessionMessages(workspaceId, agentId, sessionId db.Id, req *BatchCreateSessionMessagesRequest) ([]*SessionMessageResponse, error) {
	session, err := m.loadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	responses := make([]*SessionMessageResponse, 0, len(req.Messages))
	var lastCreatedAt time.Time
	for _, item := range req.Messages {
		message := &db.SessionMessage{
			SessionId: session.Id,
			Role:      item.Role,
			Content:   item.Content,
			Meta:      item.Meta,
		}

		err = m.MessageRepository.Create(message)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedCreateSessionMessage, err)
		}

		lastCreatedAt = message.CreatedAt
		responses = append(responses, &SessionMessageResponse{
			Id:        message.Id,
			SessionId: message.SessionId,
			Role:      message.Role,
			Content:   message.Content,
			Meta:      util.JSONRawFromString(message.Meta),
			CreatedAt: message.CreatedAt.UTC().Format(time.RFC3339),
		})

		err = m.EnqueueIndexTask(workspaceId, message.Id)
		if err != nil {
			return nil, err
		}

		if message.Role == "user" {
			err = m.EnqueueGenerateMemoryTask(workspaceId, message.Id)
			if err != nil {
				return nil, err
			}
		}
	}

	err = m.SessionRepository.TouchLastActivityAt(session.Id, lastCreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateSessionMessage, err)
	}

	return responses, nil
}

// GetSessionMessage returns one message by id.
func (m *Message) GetSessionMessage(workspaceId, agentId, sessionId, messageId db.Id) (*SessionMessageResponse, error) {
	_, err := m.loadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	message, err := m.MessageRepository.GetById(messageId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetSessionMessage, err)
	}
	if message == nil || message.SessionId != sessionId {
		return nil, ErrSessionMessageNotFound
	}

	return &SessionMessageResponse{
		Id:        message.Id,
		SessionId: message.SessionId,
		Role:      message.Role,
		Content:   message.Content,
		Meta:      util.JSONRawFromString(message.Meta),
		CreatedAt: message.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// UpdateSessionMessage updates a message by id.
func (m *Message) UpdateSessionMessage(workspaceId, agentId, sessionId, messageId db.Id, req *UpdateSessionMessageRequest) (*SessionMessageResponse, error) {
	_, err := m.loadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	message, err := m.MessageRepository.GetById(messageId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdateSessionMessage, err)
	}
	if message == nil || message.SessionId != sessionId {
		return nil, ErrSessionMessageNotFound
	}

	if lo.IsNotEmpty(req.Role) {
		message.Role = req.Role
	}
	message.Content = req.Content
	message.Meta = req.Meta

	err = m.MessageRepository.Update(message)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdateSessionMessage, err)
	}

	return &SessionMessageResponse{
		Id:        message.Id,
		SessionId: message.SessionId,
		Role:      message.Role,
		Content:   message.Content,
		Meta:      util.JSONRawFromString(message.Meta),
		CreatedAt: message.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// DeleteSessionMessage deletes a message by id.
func (m *Message) DeleteSessionMessage(workspaceId, agentId, sessionId, messageId db.Id) error {
	_, err := m.loadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return err
	}

	message, err := m.MessageRepository.GetById(messageId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteSessionMessage, err)
	}
	if message == nil || message.SessionId != sessionId {
		return ErrSessionMessageNotFound
	}

	err = m.MessageRepository.Delete(messageId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteSessionMessage, err)
	}

	return m.EnqueueDeleteTask(workspaceId, messageId)
}

// BatchDeleteSessionMessages deletes multiple messages by id.
func (m *Message) BatchDeleteSessionMessages(workspaceId, agentId, sessionId db.Id, req *BatchDeleteSessionMessagesRequest) error {
	_, err := m.loadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return err
	}

	for _, id := range req.IDs {
		message, err := m.MessageRepository.GetById(db.Id(id))
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFailedDeleteSessionMessage, err)
		}
		if message == nil || message.SessionId != sessionId {
			return ErrSessionMessageNotFound
		}

		err = m.MessageRepository.Delete(db.Id(id))
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFailedDeleteSessionMessage, err)
		}

		err = m.EnqueueDeleteTask(workspaceId, db.Id(id))
		if err != nil {
			return err
		}
	}

	return nil
}

// ListSessionMessages returns paginated messages for a session.
func (m *Message) ListSessionMessages(workspaceId, agentId, sessionId db.Id, limit, offset int) (*ListSessionMessagesResponse, error) {
	_, err := m.loadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	total, err := m.MessageRepository.CountBySessionId(sessionId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListSessionMessages, err)
	}

	messages, err := m.MessageRepository.ListBySessionId(sessionId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListSessionMessages, err)
	}

	list := make([]*SessionMessageResponse, 0, len(messages))
	for _, item := range messages {
		list = append(list, &SessionMessageResponse{
			Id:        item.Id,
			SessionId: item.SessionId,
			Role:      item.Role,
			Content:   item.Content,
			Meta:      util.JSONRawFromString(item.Meta),
			CreatedAt: item.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	return &ListSessionMessagesResponse{Messages: list, Total: total}, nil
}

// SearchSessionMessages finds relevant session messages for a semantic query.
func (m *Message) SearchSessionMessages(ctx context.Context, workspaceId, agentId, sessionId db.Id, query string, limit int) (*SearchSessionMessagesResponse, error) {
	session, err := m.loadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	hits, usage, err := m.Agent.SearchSessionMessages(ctx, agent.SearchSessionRequest{
		WorkspaceId: workspaceId,
		AgentId:     agentId,
		SessionId:   session.Id,
		Query:       query,
		Limit:       uint64(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedSearchSessionMessages, err)
	}

	err = NewUsage().IncrementAIUsage(
		m.UsageRepository,
		m.SubscriptionRepository,
		workspaceId,
		usage.TotalTokens,
		usage.Cost,
	)
	if err != nil {
		log.Warn().
			Err(err).
			Str("workspaceId", workspaceId.String()).
			Int64("total_tokens", usage.TotalTokens).
			Int64("total_cost", usage.Cost).
			Msg("Failed to record AI usage for session message search")
	}

	results := make([]SearchSessionMessageHit, 0, len(hits))
	for _, hit := range hits {
		results = append(results, SearchSessionMessageHit{
			Id:    hit.MessageId,
			Role:  hit.Role,
			Score: hit.Score,
			Text:  hit.Text,
		})
	}

	return &SearchSessionMessagesResponse{
		Results: results,
	}, nil
}

// EnqueueIndexTask enqueues an async task to index a session message.
func (m *Message) EnqueueIndexTask(workspaceId, messageId db.Id) error {
	_, err := GetAsyncMgr().CreateTask(&CreateAsyncTaskOptions{
		WorkspaceId: workspaceId,
		Type:        "task.agent.session.message.index",
		Priority:    AsyncPriorityLow,
		Payload: map[string]string{
			"messageId": messageId.String(),
		},
	})

	return err
}

// EnqueueDeleteTask enqueues an async task to delete session message vectors.
func (m *Message) EnqueueDeleteTask(workspaceId, messageId db.Id) error {
	_, err := GetAsyncMgr().CreateTask(&CreateAsyncTaskOptions{
		WorkspaceId: workspaceId,
		Type:        "task.agent.session.message.delete",
		Priority:    AsyncPriorityLow,
		Payload: map[string]string{
			"messageId": messageId.String(),
		},
	})

	return err
}

// EnqueueGenerateMemoryTask enqueues an async task to generate memory from a message.
func (m *Message) EnqueueGenerateMemoryTask(workspaceId, messageId db.Id) error {
	_, err := GetAsyncMgr().CreateTask(&CreateAsyncTaskOptions{
		WorkspaceId: workspaceId,
		Type:        "task.agent.session.memory.generate",
		Priority:    AsyncPriorityHigh,
		Payload: map[string]string{
			"messageId": messageId.String(),
		},
	})

	return err
}
