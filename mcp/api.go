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

// NewAPIClient builds a client pointed at the management server.
func NewAPIClient(mgt string, timeout time.Duration) *APIClient {
	return &APIClient{
		baseURL: strings.TrimRight(mgt, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// GetSessionByExternalID gets a session by external id.
func (c *APIClient) GetSessionByExternalID(ctx context.Context, scope Scope) (*Session, int, error) {
	path := fmt.Sprintf(
		"/api/v1/workspaces/%s/agents/%s/sessions/by-labels?id=%s",
		scope.WorkspaceID,
		scope.AgentID,
		url.QueryEscape(scope.SessionID),
	)

	var session Session
	status, err := c.Request(ctx, http.MethodGet, path, scope.AccessKey, nil, &session)
	if err != nil {
		return nil, 0, err
	}

	return &session, status, nil
}

// CreateSession creates a session with the scoped external id.
func (c *APIClient) CreateSession(ctx context.Context, scope Scope) (*Session, int, error) {
	path := fmt.Sprintf(
		"/api/v1/workspaces/%s/agents/%s/sessions?id=%s",
		scope.WorkspaceID,
		scope.AgentID,
		url.QueryEscape(scope.SessionID),
	)

	var session Session
	status, err := c.Request(ctx, http.MethodPost, path, scope.AccessKey, nil, &session)
	if err != nil {
		return nil, 0, err
	}

	return &session, status, nil
}

// DeleteSession deletes the scoped session (and its memories).
func (c *APIClient) DeleteSession(ctx context.Context, scope Scope) error {
	path := fmt.Sprintf(
		"/api/v1/workspaces/%s/agents/%s/sessions/by-labels?id=%s",
		scope.WorkspaceID,
		scope.AgentID,
		url.QueryEscape(scope.SessionID),
	)

	status, err := c.Request(ctx, http.MethodDelete, path, scope.AccessKey, nil, nil)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("delete session: unexpected status %d", status)
	}

	return nil
}

// EnsureSession returns the internal session id, creating the session when missing.
func (c *APIClient) EnsureSession(ctx context.Context, scope Scope) (string, error) {
	session, status, err := c.GetSessionByExternalID(ctx, scope)
	if err != nil {
		return "", err
	}
	if status == http.StatusOK {
		return session.Id, nil
	}

	session, status, err = c.CreateSession(ctx, scope)
	if err != nil {
		return "", err
	}
	if status == http.StatusCreated {
		return session.Id, nil
	}
	if status == http.StatusConflict {
		session, status, err = c.GetSessionByExternalID(ctx, scope)
		if err != nil {
			return "", err
		}
		if status == http.StatusOK {
			return session.Id, nil
		}
	}

	return "", fmt.Errorf("ensure session: unexpected status %d", status)
}

// ResolveSession returns the internal session id for Scope.SessionID.
func (c *APIClient) ResolveSession(ctx context.Context, scope Scope) (string, error) {
	session, status, err := c.GetSessionByExternalID(ctx, scope)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("resolve session: unexpected status %d", status)
	}

	return session.Id, nil
}

// CreateMemory stores a session memory.
func (c *APIClient) CreateMemory(ctx context.Context, scope Scope, kind, content string, meta *string) (*Memory, error) {
	sessionID, err := c.EnsureSession(ctx, scope)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(MemoryBody{
		Kind:    kind,
		Content: content,
		Meta:    meta,
	})
	if err != nil {
		return nil, err
	}

	var memory Memory
	status, err := c.Request(ctx, http.MethodPost, c.MemoriesPath(scope, sessionID), scope.AccessKey, body, &memory)
	if err != nil {
		return nil, err
	}
	if status != http.StatusCreated {
		return nil, fmt.Errorf("create memory: unexpected status %d", status)
	}

	return &memory, nil
}

// ListMemories lists session memories.
func (c *APIClient) ListMemories(ctx context.Context, scope Scope, limit, offset int) (*ListMemories, error) {
	sessionID, err := c.ResolveSession(ctx, scope)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s?limit=%d&offset=%d", c.MemoriesPath(scope, sessionID), limit, offset)

	var list ListMemories
	status, err := c.Request(ctx, http.MethodGet, path, scope.AccessKey, nil, &list)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("list memories: unexpected status %d", status)
	}

	return &list, nil
}

// GetMemory retrieves one session memory.
func (c *APIClient) GetMemory(ctx context.Context, scope Scope, memoryID string) (*Memory, error) {
	sessionID, err := c.ResolveSession(ctx, scope)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", c.MemoriesPath(scope, sessionID), memoryID)

	var memory Memory
	status, err := c.Request(ctx, http.MethodGet, path, scope.AccessKey, nil, &memory)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("get memory: unexpected status %d", status)
	}

	return &memory, nil
}

// UpdateMemory updates a session memory.
func (c *APIClient) UpdateMemory(ctx context.Context, scope Scope, memoryID, kind, content string, meta *string) (*Memory, error) {
	sessionID, err := c.ResolveSession(ctx, scope)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(MemoryBody{
		Kind:    kind,
		Content: content,
		Meta:    meta,
	})
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", c.MemoriesPath(scope, sessionID), memoryID)

	var memory Memory
	status, err := c.Request(ctx, http.MethodPut, path, scope.AccessKey, body, &memory)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("update memory: unexpected status %d", status)
	}

	return &memory, nil
}

// SearchMemories searches session memories.
func (c *APIClient) SearchMemories(ctx context.Context, scope Scope, query string, limit int) (*SearchMemories, error) {
	sessionID, err := c.ResolveSession(ctx, scope)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(SearchBody{
		Query: query,
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	path := c.MemoriesPath(scope, sessionID) + "/search"

	var result SearchMemories
	status, err := c.Request(ctx, http.MethodPost, path, scope.AccessKey, body, &result)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("search memories: unexpected status %d", status)
	}

	return &result, nil
}

// DeleteMemory deletes a session memory.
func (c *APIClient) DeleteMemory(ctx context.Context, scope Scope, memoryID string) error {
	sessionID, err := c.ResolveSession(ctx, scope)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", c.MemoriesPath(scope, sessionID), memoryID)

	status, err := c.Request(ctx, http.MethodDelete, path, scope.AccessKey, nil, nil)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("delete memory: unexpected status %d", status)
	}

	return nil
}

// MemoriesPath builds the memories collection path for a session.
func (c *APIClient) MemoriesPath(scope Scope, sessionID string) string {
	return fmt.Sprintf(
		"/api/v1/workspaces/%s/agents/%s/sessions/%s/memories",
		scope.WorkspaceID,
		scope.AgentID,
		sessionID,
	)
}

// Request performs an HTTP request against the management server.
func (c *APIClient) Request(ctx context.Context, method, path, accessKey string, body []byte, out any) (int, error) {
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

// EncodeMeta marshals optional metadata for memory create/update bodies.
func EncodeMeta(metadata map[string]any) (*string, error) {
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
