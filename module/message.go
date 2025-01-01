// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/util"
)

var (
	ErrSessionMessageNotFound     = errors.New("session message not found")
	ErrFailedCreateSessionMessage = errors.New("failed create session message")
	ErrFailedListSessionMessages  = errors.New("failed list session messages")
	ErrFailedGetSessionMessage    = errors.New("failed get session message")
	ErrFailedUpdateSessionMessage = errors.New("failed update session message")
	ErrFailedDeleteSessionMessage = errors.New("failed delete session message")
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
	MessageRepository db.SessionMessageRepository
	AgentSession      *AgentSession
}

// NewMessage creates a session message module with the given repositories.
func NewMessage(mr db.SessionMessageRepository, as *AgentSession) *Message {
	return &Message{
		MessageRepository: mr,
		AgentSession:      as,
	}
}

// CreateSessionMessage creates a message in a session.
func (m *Message) CreateSessionMessage(workspaceId, agentId, sessionId db.Id, req *CreateSessionMessageRequest) (*SessionMessageResponse, error) {
	session, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
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

	return sessionMessageResponse(message), nil
}

// BatchCreateSessionMessages creates multiple messages in a session.
func (m *Message) BatchCreateSessionMessages(workspaceId, agentId, sessionId db.Id, req *BatchCreateSessionMessagesRequest) ([]*SessionMessageResponse, error) {
	session, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
	if err != nil {
		return nil, err
	}

	responses := make([]*SessionMessageResponse, 0, len(req.Messages))
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

		responses = append(responses, sessionMessageResponse(message))
	}

	return responses, nil
}

// GetSessionMessage returns one message by id.
func (m *Message) GetSessionMessage(workspaceId, agentId, sessionId, messageId db.Id) (*SessionMessageResponse, error) {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
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

	return sessionMessageResponse(message), nil
}

// UpdateSessionMessage updates a message by id.
func (m *Message) UpdateSessionMessage(workspaceId, agentId, sessionId, messageId db.Id, req *UpdateSessionMessageRequest) (*SessionMessageResponse, error) {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
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

	if req.Role != "" {
		message.Role = req.Role
	}
	message.Content = req.Content
	message.Meta = req.Meta

	err = m.MessageRepository.Update(message)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdateSessionMessage, err)
	}

	return sessionMessageResponse(message), nil
}

// DeleteSessionMessage deletes a message by id.
func (m *Message) DeleteSessionMessage(workspaceId, agentId, sessionId, messageId db.Id) error {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
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

	_, err = GetAsyncMgr().CreateTask(&CreateAsyncTaskOptions{
		WorkspaceId: workspaceId,
		Type:        "task.agent.session.message.delete",
		Priority:    AsyncPriorityLow,
		Payload: map[string]string{
			"messageId": messageId.String(),
		},
	})

	return err
}

// BatchDeleteSessionMessages deletes multiple messages by id.
func (m *Message) BatchDeleteSessionMessages(workspaceId, agentId, sessionId db.Id, req *BatchDeleteSessionMessagesRequest) error {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
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

		_, err = GetAsyncMgr().CreateTask(&CreateAsyncTaskOptions{
			WorkspaceId: workspaceId,
			Type:        "task.agent.session.message.delete",
			Priority:    AsyncPriorityLow,
			Payload: map[string]string{
				"messageId": id,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// ListSessionMessages returns paginated messages for a session.
func (m *Message) ListSessionMessages(workspaceId, agentId, sessionId db.Id, limit, offset int) (*ListSessionMessagesResponse, error) {
	_, err := m.AgentSession.LoadAgentSession(workspaceId, agentId, sessionId)
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
		list = append(list, sessionMessageResponse(item))
	}

	return &ListSessionMessagesResponse{Messages: list, Total: total}, nil
}

func sessionMessageResponse(message *db.SessionMessage) *SessionMessageResponse {
	return &SessionMessageResponse{
		Id:        message.Id,
		SessionId: message.SessionId,
		Role:      message.Role,
		Content:   message.Content,
		Meta:      util.JSONRawFromString(message.Meta),
		CreatedAt: message.CreatedAt.UTC().Format(time.RFC3339),
	}
}
