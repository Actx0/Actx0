// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// APIClient calls the actx0 management HTTP API.
type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

type sessionDTO struct {
	ID string `json:"id"`
}

type memoryDTO struct {
	ID        string          `json:"id"`
	SessionID string          `json:"sessionId"`
	Kind      string          `json:"kind"`
	Content   string          `json:"content"`
	Meta      json.RawMessage `json:"meta,omitempty"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

type listMemoriesDTO struct {
	Memories []memoryDTO `json:"memories"`
	Meta     struct {
		Limit  int   `json:"limit"`
		Offset int   `json:"offset"`
		Total  int64 `json:"total"`
	} `json:"_meta"`
}

type searchHitDTO struct {
	ID    string  `json:"id"`
	Kind  string  `json:"kind"`
	Score float32 `json:"score"`
	Text  string  `json:"text"`
}

type searchMemoriesDTO struct {
	Results []searchHitDTO `json:"results"`
}

type memoryBody struct {
	Kind    string  `json:"kind,omitempty"`
	Content string  `json:"content"`
	Meta    *string `json:"meta,omitempty"`
}

type searchBody struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

// NewAPIClient builds a client pointed at the management server.
func NewAPIClient(mgt string, timeout time.Duration) *APIClient {
	return &APIClient{
		baseURL: strings.TrimRight(mgt, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// EnsureSession returns the internal session id, creating the session when missing.
// Scope.SessionID is treated as the caller external session id.
func (c *APIClient) EnsureSession(ctx context.Context, scope Scope) (string, error) {
	session, status, err := c.getSessionByExternalID(ctx, scope)
	if err != nil {
		return "", err
	}
	if status == http.StatusOK {
		return session.ID, nil
	}

	session, status, err = c.createSession(ctx, scope)
	if err != nil {
		return "", err
	}
	if status == http.StatusCreated {
		return session.ID, nil
	}
	if status == http.StatusConflict {
		session, status, err = c.getSessionByExternalID(ctx, scope)
		if err != nil {
			return "", err
		}
		if status == http.StatusOK {
			return session.ID, nil
		}
	}

	return "", fmt.Errorf("ensure session: unexpected status %d", status)
}

// ResolveSession returns the internal session id for Scope.SessionID.
func (c *APIClient) ResolveSession(ctx context.Context, scope Scope) (string, error) {
	session, status, err := c.getSessionByExternalID(ctx, scope)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("resolve session: unexpected status %d", status)
	}

	return session.ID, nil
}

// CreateMemory stores a session memory via the actx0 API.
func (c *APIClient) CreateMemory(ctx context.Context, scope Scope, kind, content string, meta *string) (*memoryDTO, error) {
	sessionID, err := c.EnsureSession(ctx, scope)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(memoryBody{
		Kind:    kind,
		Content: content,
		Meta:    meta,
	})
	if err != nil {
		return nil, err
	}

	var memory memoryDTO
	status, err := c.doJSON(ctx, http.MethodPost, c.memoriesPath(scope, sessionID), scope.AccessKey, body, &memory)
	if err != nil {
		return nil, err
	}
	if status != http.StatusCreated {
		return nil, fmt.Errorf("create memory: unexpected status %d", status)
	}

	return &memory, nil
}

// ListMemories lists session memories via the actx0 API.
func (c *APIClient) ListMemories(ctx context.Context, scope Scope, limit, offset int) (*listMemoriesDTO, error) {
	sessionID, err := c.ResolveSession(ctx, scope)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s?limit=%d&offset=%d", c.memoriesPath(scope, sessionID), limit, offset)

	var list listMemoriesDTO
	status, err := c.doJSON(ctx, http.MethodGet, path, scope.AccessKey, nil, &list)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("list memories: unexpected status %d", status)
	}

	return &list, nil
}

// GetMemory retrieves one session memory via the actx0 API.
func (c *APIClient) GetMemory(ctx context.Context, scope Scope, memoryID string) (*memoryDTO, error) {
	sessionID, err := c.ResolveSession(ctx, scope)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", c.memoriesPath(scope, sessionID), memoryID)

	var memory memoryDTO
	status, err := c.doJSON(ctx, http.MethodGet, path, scope.AccessKey, nil, &memory)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("get memory: unexpected status %d", status)
	}

	return &memory, nil
}

// UpdateMemory updates a session memory via the actx0 API.
func (c *APIClient) UpdateMemory(ctx context.Context, scope Scope, memoryID, kind, content string, meta *string) (*memoryDTO, error) {
	sessionID, err := c.ResolveSession(ctx, scope)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(memoryBody{
		Kind:    kind,
		Content: content,
		Meta:    meta,
	})
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", c.memoriesPath(scope, sessionID), memoryID)

	var memory memoryDTO
	status, err := c.doJSON(ctx, http.MethodPut, path, scope.AccessKey, body, &memory)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("update memory: unexpected status %d", status)
	}

	return &memory, nil
}

// SearchMemories searches session memories via the actx0 API.
func (c *APIClient) SearchMemories(ctx context.Context, scope Scope, query string, limit int) (*searchMemoriesDTO, error) {
	sessionID, err := c.ResolveSession(ctx, scope)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(searchBody{
		Query: query,
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	path := c.memoriesPath(scope, sessionID) + "/search"

	var result searchMemoriesDTO
	status, err := c.doJSON(ctx, http.MethodPost, path, scope.AccessKey, body, &result)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("search memories: unexpected status %d", status)
	}

	return &result, nil
}

// DeleteMemory deletes a session memory via the actx0 API.
func (c *APIClient) DeleteMemory(ctx context.Context, scope Scope, memoryID string) error {
	sessionID, err := c.ResolveSession(ctx, scope)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", c.memoriesPath(scope, sessionID), memoryID)

	status, err := c.doJSON(ctx, http.MethodDelete, path, scope.AccessKey, nil, nil)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("delete memory: unexpected status %d", status)
	}

	return nil
}

// DeleteSession deletes the scoped session (and its memories) via the actx0 API.
func (c *APIClient) DeleteSession(ctx context.Context, scope Scope) error {
	path := fmt.Sprintf(
		"/api/v1/workspaces/%s/agents/%s/sessions/by-labels?id=%s",
		scope.WorkspaceID,
		scope.AgentID,
		url.QueryEscape(scope.SessionID),
	)

	status, err := c.doJSON(ctx, http.MethodDelete, path, scope.AccessKey, nil, nil)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("delete session: unexpected status %d", status)
	}

	return nil
}

func (c *APIClient) memoriesPath(scope Scope, sessionID string) string {
	return fmt.Sprintf(
		"/api/v1/workspaces/%s/agents/%s/sessions/%s/memories",
		scope.WorkspaceID,
		scope.AgentID,
		sessionID,
	)
}

func (c *APIClient) getSessionByExternalID(ctx context.Context, scope Scope) (*sessionDTO, int, error) {
	path := fmt.Sprintf(
		"/api/v1/workspaces/%s/agents/%s/sessions/by-labels?id=%s",
		scope.WorkspaceID,
		scope.AgentID,
		url.QueryEscape(scope.SessionID),
	)

	var session sessionDTO
	status, err := c.doJSON(ctx, http.MethodGet, path, scope.AccessKey, nil, &session)
	if err != nil {
		return nil, 0, err
	}

	return &session, status, nil
}

func (c *APIClient) createSession(ctx context.Context, scope Scope) (*sessionDTO, int, error) {
	path := fmt.Sprintf(
		"/api/v1/workspaces/%s/agents/%s/sessions?id=%s",
		scope.WorkspaceID,
		scope.AgentID,
		url.QueryEscape(scope.SessionID),
	)

	var session sessionDTO
	status, err := c.doJSON(ctx, http.MethodPost, path, scope.AccessKey, nil, &session)
	if err != nil {
		return nil, 0, err
	}

	return &session, status, nil
}

func (c *APIClient) doJSON(
	ctx context.Context,
	method, path, accessKey string,
	body []byte,
	out any,
) (int, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return 0, err
	}

	req.Header.Set(HeaderAccessKey, accessKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if out == nil || resp.StatusCode == http.StatusNoContent {
		io.Copy(io.Discard, resp.Body)
		return resp.StatusCode, nil
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		err = json.NewDecoder(resp.Body).Decode(out)
		return resp.StatusCode, err
	}

	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}

func encodeMeta(metadata map[string]any) (*string, error) {
	if len(metadata) == 0 {
		return nil, nil
	}

	raw, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	s := string(raw)
	return &s, nil
}

func memoryRecordFromDTO(item memoryDTO) MemoryRecord {
	record := MemoryRecord{
		MemoryID:  item.ID,
		Kind:      item.Kind,
		Memory:    item.Content,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}

	if len(item.Meta) > 0 {
		var metadata map[string]any
		_ = json.Unmarshal(item.Meta, &metadata)
		record.Metadata = metadata
	}

	return record
}
